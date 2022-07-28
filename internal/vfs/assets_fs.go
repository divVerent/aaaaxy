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

package vfs

import (
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
)

var (
	dumpEmbeddedAssets         = flag.String("dump_embedded_assets", "", "dump all embedded assets and license information to the given directory instead of running the game")
	cheatReplaceEmbeddedAssets = flag.String("cheat_replace_embedded_assets", "", "if set, embedded assets are skipped and this directory is used as assets root instead")
)

type fsRoot struct {
	name     string
	filesys  fs.FS
	root     string
	toPrefix string
}

func (f fsRoot) String() string {
	return fmt.Sprintf("%v(%v->%v)", f.name, f.root, f.toPrefix)
}

var (
	assetDirs []fsRoot
)

func dumpAssetsFrom(dir fsRoot) error {
	return fs.WalkDir(dir.filesys, dir.root, func(p string, d fs.DirEntry, err error) error {
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
		fIn, err := dir.filesys.Open(p)
		if err != nil {
			return err
		}
		defer fIn.Close()
		out := path.Join(*dumpEmbeddedAssets, dir.toPrefix+relPath)
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
	for _, dir := range assetDirs {
		err := dumpAssetsFrom(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// initAssets initializes the VFS. Must run after loading the assets.
func initAssets() error {
	if *cheatReplaceEmbeddedAssets == "" {
		var err error
		assetDirs, err = initAssetsFS()
		if err != nil {
			return err
		}
	} else {
		assetDirs = []fsRoot{
			{
				name:     "replace",
				filesys:  os.DirFS(*cheatReplaceEmbeddedAssets),
				root:     "/",
				toPrefix: "",
			},
		}
	}

	log.Infof("asset search path: %v", assetDirs)

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
	var err error
	for _, dir := range assetDirs {
		if !strings.HasPrefix(vfsPath, dir.toPrefix) {
			continue
		}
		relPath := strings.TrimPrefix(vfsPath, dir.toPrefix)
		var f fs.File
		f, err = dir.filesys.Open(path.Join(dir.root, relPath))
		if err != nil {
			continue
		}
		rsc, ok := f.(ReadSeekCloser)
		if ok {
			return rsc, nil
		}
		info, err := f.Stat()
		if err == nil && info.IsDir() {
			return nil, fmt.Errorf("could not open %v: is a directory", vfsPath)
		}
		return nil, fmt.Errorf("could not open %v: internal error (go:embed doesn't yield ReadSeekCloser)", vfsPath)
	}
	return nil, fmt.Errorf("could not open %v: %w", vfsPath, err)
}

// readDir lists all files in a directory. Returns their VFS names, NOT full paths!
func readDir(vfsPath string) ([]string, error) {
	var results []string
	for _, dir := range assetDirs {
		if !strings.HasPrefix(vfsPath, dir.toPrefix) {
			continue
		}
		relPath := strings.TrimPrefix(vfsPath, dir.toPrefix)
		content, err := fs.ReadDir(dir.filesys, path.Join(dir.root, relPath))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("could not scan %v in %v: %v", vfsPath, dir, err)
			}
			continue
		}
		for _, info := range content {
			results = append(results, dir.toPrefix+info.Name())
		}
	}
	sort.Strings(results)
	return results, nil
}
