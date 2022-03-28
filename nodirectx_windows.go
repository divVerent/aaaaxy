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

//go:build windows
// +build windows

package main

import (
	"os"
)

func init() {
	// Turn off DirectX for now, it's not stable in Ebiten yet.
	// This can be overridden by explicitly setting EBITEN_GRAPHICS_LIBRARY to "auto" or "directx".
	// TODO(rpolzer): remove when Ebiten 2.3 is released.
	driver := os.LookupEnv("EBITEN_GRAPHICS_LIBRARY")
	if driver == "" {
		os.SetEnv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	}
}
