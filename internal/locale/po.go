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

// Type is the type of a locale translator.
type Type = *gotext.Po

// IType is the type of a locale translator interface.
type IType = gotext.Translator

// G is the translation of the game.
var G Type

// GI is the translation of the game but as an interface.
// Used to work around go vet as it can't figure out
// that G.Get isn't a printf-like function if it gets only one arg.
// See https://github.com/golang/go/issues/57288
var GI IType

// L is the translation of the levels.
var L Type

// Active is the name of the current language.
var Active Lingua

// ResetLanguage makes everything English again.
func ResetLanguage() {
	G = gotext.NewPo()
	GI = G
	L = gotext.NewPo()
	Active = ""
}

func init() {
	// The translators must always be valid objects.
	ResetLanguage()
}
