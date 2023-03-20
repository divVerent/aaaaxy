// Copyright 2022 Google LLC
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

//go:build zip && (darwin || ios)
// +build zip
// +build darwin ios

package vfs

import (
	"errors"
	"os"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/NSBundle.h>
#import <Foundation/NSString.h>

#include <stdlib.h>
#include <string.h>

const char *assets_zip_path() {
	NSString *path = [[NSBundle mainBundle] pathForResource: @"aaaaxy" ofType: @"dat"];
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

func openAssetsZip() (*os.File, error) {
	pathCStr := C.assets_zip_path()
	if pathCStr == nil {
		return nil, errors.New("could not find aaaaxy.dat: PathForResource:ofType: failed")
	}
	defer C.free(unsafe.Pointer(pathCStr))
	return os.Open(C.GoString(pathCStr))
}
