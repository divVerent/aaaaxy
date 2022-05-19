// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build embed
// +build embed

package vfs

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"

	"github.com/divVerent/aaaaxy/assets"
	"github.com/divVerent/aaaaxy/licenses"
	"github.com/divVerent/aaaaxy/third_party"
)

var (
	dumpEmbeddedAssets         = flag.String("dump_embedded_assets", "", "dump all embedded assets and license information to the given directory instead of running the game")
	cheatReplaceEmbeddedAssets = flag.String("cheat_replace_embedded_assets", "", "if set, embedded assets are skipped and this directory is used as assets root instead")
)

type fsRoot struct {
	fs   *embed.FS
	root string
}

var (
	embeddedAssetDirs []fsRoot
)

func dumpAssetsFrom(dir fsRoot) error {
	return fs.WalkDir(dir.fs, dir.root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath := p
		if strings.HasPrefix(p, dir.root+"/") {
			relPath = p[len(dir.root)+1:]
		}
		fIn, err := dir.fs.Open(p)
		if err != nil {
			return err
		}
		defer fIn.Close()
		out := path.Join(*dumpEmbeddedAssets, relPath)
		log.Infof("%v => %v", p, out)
		err = os.MkdirAll(path.Dir(out), 0777)
		if err != nil {
			return err
		}
		fOut, err := os.Create(out)
		if err != nil {
			return err
		}
		_, err = io.Copy(fOut, fIn)
		if err != nil {
			fOut.Close()
			return err
		}
		return fOut.Close()
	})
}

func dumpAssets() error {
	for _, dir := range embeddedAssetDirs {
		err := dumpAssetsFrom(dir)
		if err != nil {
			return err
		}
	}
	return dumpAssetsFrom(fsRoot{
		fs:   &licenses.FS,
		root: ".",
	})
}

// initAssets initializes the VFS. Must run after loading the assets.
func initAssets() error {
	if *cheatReplaceEmbeddedAssets != "" {
		return initLocalAssets([]string{*cheatReplaceEmbeddedAssets})
	}

	embeddedAssetDirs = []fsRoot{{
		fs:   &assets.FS,
		root: ".",
	}}
	content, err := third_party.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("could not find embedded third party directory: %v", err)
	}
	roots := []string{}
	for _, info := range content {
		if !info.IsDir() {
			continue
		}
		root := path.Join(info.Name(), "assets")
		embeddedAssetDirs = append(embeddedAssetDirs, fsRoot{
			fs:   &third_party.FS,
			root: root,
		})
		roots = append(roots, root)
	}
	log.Debugf("embedded asset search path: %v", roots)

	if *dumpEmbeddedAssets != "" {
		err := dumpAssets()
		if err != nil {
			return err
		}
		log.Infof("requested an asset dump - not running the game")
		return exitstatus.RegularTermination
	}

	return nil
}

// load loads a file from the VFS.
func load(vfsPath string) (ReadSeekCloser, error) {
	if *cheatReplaceEmbeddedAssets != "" {
		return loadLocal(vfsPath)
	}

	var err error
	for _, dir := range embeddedAssetDirs {
		var f fs.File
		f, err = dir.fs.Open(path.Join(dir.root, vfsPath))
		if err != nil {
			continue
		}
		rsc, ok := f.(ReadSeekCloser)
		if !ok {
			info, err := f.Stat()
			if err == nil && info.IsDir() {
				return nil, fmt.Errorf("could not open embed:%v: is a directory", vfsPath)
			}
			return nil, fmt.Errorf("could not open embed:%v: internal error (go:embed doesn't yield ReadSeekCloser)", vfsPath)
		}
		return rsc, nil
	}
	return nil, fmt.Errorf("could not open embed:%v: %w", vfsPath, err)
}

// readDir lists all files in a directory. Returns their VFS names, NOT full paths!
func readDir(vfsPath string) ([]string, error) {
	if *cheatReplaceEmbeddedAssets != "" {
		return readLocalDir(vfsPath)
	}

	var results []string
	for _, dir := range embeddedAssetDirs {
		content, err := dir.fs.ReadDir(path.Join(dir.root, vfsPath))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("could not scan embed:%v:%v: %v", vfsPath, dir, err)
			}
			continue
		}
		for _, info := range content {
			results = append(results, info.Name())
		}
	}
	sort.Strings(results)
	return results, nil
}
