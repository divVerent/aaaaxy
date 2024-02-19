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

//go:build zip && !embed && android
// +build zip,!embed,android

package vfs

// Following code is copied from
// golang.org/x/mobile/asset/asset_android.go
// but changed to work with ebitenmobile.

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
#cgo LDFLAGS: -landroid
#include <android/asset_manager.h>
#include <android/asset_manager_jni.h>
#include <jni.h>
#include <stdlib.h>

static AAssetManager* asset_manager_init(uintptr_t java_vm, uintptr_t jni_env, jobject ctx) {
	JavaVM* vm = (JavaVM*)java_vm;
	JNIEnv* env = (JNIEnv*)jni_env;

	// Equivalent to:
	//	assetManager = ctx.getResources().getAssets();
	jclass ctx_clazz = (*env)->FindClass(env, "android/content/Context");
	jmethodID getres_id = (*env)->GetMethodID(env, ctx_clazz, "getResources", "()Landroid/content/res/Resources;");
	jobject res = (*env)->CallObjectMethod(env, ctx, getres_id);
	jclass res_clazz = (*env)->FindClass(env, "android/content/res/Resources");
	jmethodID getam_id = (*env)->GetMethodID(env, res_clazz, "getAssets", "()Landroid/content/res/AssetManager;");
	jobject am = (*env)->CallObjectMethod(env, res, getam_id);

	// Pin the AssetManager and load an AAssetManager from it.
	am = (*env)->NewGlobalRef(env, am);
	return AAssetManager_fromJava(env, am);
}
*/
import "C"

import (
	"errors"
	"io"
	"io/fs"
	"sync"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/divVerent/aaaaxy/internal/log"
)

var assetOnce sync.Once

// asset_manager is the asset manager of the app.
var assetManager *C.AAssetManager

func assetInit() {
	err := mobile.RunOnJVM(func(vm, env, ctx uintptr) error {
		assetManager = C.asset_manager_init(C.uintptr_t(vm), C.uintptr_t(env), C.jobject(ctx))
		return nil
	})
	if err != nil {
		log.Fatalf("asset: %v", err)
	}
}

type fileInfo int64

func (i fileInfo) Name() string       { return "asset" }
func (i fileInfo) Size() int64        { return int64(i) }
func (i fileInfo) Mode() fs.FileMode  { return 0555 }
func (i fileInfo) ModTime() time.Time { return time.Time{} }
func (i fileInfo) IsDir() bool        { return false }
func (i fileInfo) Sys() interface{}   { return nil }

type reader struct {
	ptr *C.AAsset
	mu  sync.Mutex
	pos int64
}

func (r *reader) readAt(p []byte, off int64) (int, error) {
	pos := int64(C.AAsset_seek(r.ptr, C.off_t(off), C.int(io.SeekStart)))
	if pos < 0 {
		return 0, errors.New("could not seek in aaaaxy.dat")
	}
	n := int(C.AAsset_read(r.ptr, unsafe.Pointer(&p[0]), C.size_t(len(p))))
	if n == 0 && len(p) > 0 {
		return 0, io.EOF
	}
	if n < 0 {
		return 0, errors.New("could not read from aaaaxy.dat")
	}
	return n, nil
}

func (r *reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, err := r.readAt(p, r.pos)
	r.pos += int64(n)
	return n, err
}

func (r *reader) ReadAt(p []byte, off int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readAt(p, off)
}

func (r *reader) Stat() (fs.FileInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	size := int64(C.AAsset_getLength(r.ptr))
	return fileInfo(size), nil
}

func (r *reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	C.AAsset_close(r.ptr)
	return nil
}

func openAssetsZip() (*reader, error) {
	assetOnce.Do(assetInit)
	cname := C.CString("aaaaxy.dat")
	defer C.free(unsafe.Pointer(cname))
	mode := C.int(C.AASSET_MODE_RANDOM)
	if *pinAssetsToRAM {
		mode = C.AASSET_MODE_STREAMING
	}
	ptr := C.AAssetManager_open(assetManager, cname, mode)
	if ptr == nil {
		return nil, errors.New("could not open aaaaxy.dat")
	}
	return &reader{ptr: ptr}, nil
}
