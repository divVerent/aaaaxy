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

package font

// charSetBase is the set of all characters we use in English, ignoring credits.
const charSetBase = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~τπö¾©éàã"

var (
	// charSet is the active character set. Used for font pinning.
	charSet []rune

	// charSetCached indicates whether charSet is already cached.
	charSetCached bool

	// charSetPos is the current caching position.
	charSetPos int
)
