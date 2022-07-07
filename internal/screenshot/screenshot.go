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

package screenshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func Write(img image.Image, name string) (err error) {
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to open image file %v: %w", name, err)
	}
	defer func() {
		errC := file.Close()
		if errC != nil && err == nil {
			err = fmt.Errorf("failed to close image file %v: %w", name, errC)
		}
	}()
	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to write to image file %v: %w", name, err)
	}
	return nil
}
