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

//go:build ios
// +build ios

package vfs

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/divVerent/aaaaxy/internal/log"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/NSPathUtilities.h>
#import <Foundation/NSString.h>

#include <stdlib.h>
#include <string.h>

const char *documents_path() {
	// TODO: support this containing more than one element.
	NSString *path = [NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES) elementAtIndex: 0];
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
	documentsPath string
)

func pathForReadRaw(kind StateKind, name string) (string, error) {
	return pathForWrite(kind, name)
}

func pathForWriteRaw(kind StateKind, name string) (string, error) {
	if documentsPath == "" {
		documentsPathCStr := C.documents_path()
		if documentsPathCStr == nil {
			return "", errors.New("could not find documents path")
		}
		defer C.free(unsafe.Pointer(documentsPathCStr))
		documentsPath = C.GoString(documentsPathCStr)
	}
	switch kind {
	case Config:
		return filepath.Join(documentsPath, "config", name)
	case SavedGames:
		return filepath.Join(documentsPath, "save", name)
	default:
		return "", fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}
