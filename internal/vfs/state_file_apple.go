// Copyright 2023 Google LLC
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

//go:build darwin || ios
// +build darwin ios

package vfs

import (
	"fmt"
	"path/filepath"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/NSPathUtilities.h>
#import <Foundation/NSString.h>

#include <stdlib.h>
#include <string.h>

const char *app_support_path() {
	NSArray<NSString *> *paths = NSSearchPathForDirectoriesInDomains(NSApplicationSupportDirectory, NSUserDomainMask, YES);
	if ([paths count] < 1) {
		return NULL;
	}
	NSString *path = [paths firstObject];
	if (path == nil) {
		return NULL;
	}
	const char *data = [path UTF8String];
	if (data == NULL) {
		return NULL;
	}
	return strdup(data);
}
*/
import "C"

var (
	appSupportPath string
)

func initAppSupportPath() {
	if appSupportPath == "" {
		appSupportPathCStr := C.app_support_path()
		if appSupportPathCStr == nil {
			panic("could not find documents path")
		}
		defer C.free(unsafe.Pointer(appSupportPathCStr))
		appSupportPath = C.GoString(appSupportPathCStr)
	}
}

func pathForReadRaw(kind StateKind, name string) ([]string, error) {
	initAppSupportPath()
	switch kind {
	case Config:
		return []string{
			filepath.Join(appSupportPath, "AAAAXY", "config", name),
			// This one matches state_file_xdg.go's for compatibility.
			filepath.Join(appSupportPath, "AAAAXY", name),
		}, nil
	case SavedGames:
		return []string{
			filepath.Join(appSupportPath, "AAAAXY", "save", name),
			// This one matches state_file_xdg.go's for compatibility.
			filepath.Join(appSupportPath, "AAAAXY", name),
		}, nil
	default:
		return nil, fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}

func pathForWriteRaw(kind StateKind, name string) (string, error) {
	initAppSupportPath()
	switch kind {
	case Config:
		return filepath.Join(appSupportPath, "AAAAXY", "config", name), nil
	case SavedGames:
		return filepath.Join(appSupportPath, "AAAAXY", "save", name), nil
	default:
		return "", fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}
