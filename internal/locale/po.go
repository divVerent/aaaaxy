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

package locale

import (
	"github.com/leonelquinteros/gotext"
)

// Type is the type of a locale.
type Type = *gotext.Po

// G is the translation of the game.
var G Type

// L is the translation of the levels.
var L Type

// Active is the name of the current language.
var Active Lingua

func init() {
	// Make T always a valid object.
	// Before translations are loaded, it does nothing.
	G = gotext.NewPo()
	L = gotext.NewPo()
	Active = ""
}
