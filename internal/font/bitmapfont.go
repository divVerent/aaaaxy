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

package font

import (
	"github.com/hajimehoshi/bitmapfont/v3"
)

func initBitmapfont() error {
	face := makeFace(bitmapfont.Face, 16)

	ByName["Small"] = face
	ByName["Regular"] = face
	ByName["Italic"] = face
	ByName["Bold"] = face
	ByName["Mono"] = face
	ByName["MonoSmall"] = face
	ByName["SmallCaps"] = face
	ByName["Centerprint"] = face
	ByName["CenterprintBig"] = face
	ByName["Menu"] = face
	ByName["MenuBig"] = face
	ByName["MenuSmall"] = face

	return nil
}
