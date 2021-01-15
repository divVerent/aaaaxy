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

// A simple shader to perform blurs.
package main

var Size float
var Step vec2
var Scale float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	acc := imageSrc0UnsafeAt(texCoord)
	for y := 1.0; y <= 6.0; y += 1.0 {
		if y <= Size {
			d := y * Step
			acc += imageSrc0At(texCoord - d)
			acc += imageSrc0At(texCoord + d)
		}
	}
	return acc * Scale
}
