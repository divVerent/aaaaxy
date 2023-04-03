// Copyright 2023 Google LLC
//
// Licensed under the Apache Livense, Version 2.0 (the "License");
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

package log

import (
	// Some package that includes mobileinit.
	// This ensures that init() from here runs later.
	_ "golang.org/x/mobile/app"
)

// Override logging to be public.
// That way, Console shows log messages from this game.
// Logs in AAAAXY simply never contain private info anyway.

// Following code is copied from
// golang.org/x/mobile/internal/mobileinit/mobileinit_ios.go
// but changed to log as public.

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"io"
	"log"
	"os"
	"unsafe"
)

/*
#include <stdlib.h>
#include <os/log.h>

os_log_t aaaaxy_create_os_log() {
	return os_log_create("org.golang.mobile", "os_log");
}

void aaaaxy_os_log_wrap(os_log_t log, const char *str) {
	os_log(log, "%{public}s", str);
}
*/
import "C"

type osWriter struct {
	w C.os_log_t
}

func (o osWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.aaaaxy_os_log_wrap(o.w, cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

func init() {
	log.SetOutput(io.MultiWriter(os.Stderr, osWriter{C.aaaaxy_create_os_log()}))
}
