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

package shader

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	debugUseShaders = flag.Bool("debug_use_shaders", true, "enable use of custom shaders")
)

type shaderPath = struct {
	Name             string
	ParamsSerialized string
}

var cache = map[shaderPath]*ebiten.Shader{}

func Load(name string, params map[string]string) (*ebiten.Shader, error) {
	if !*debugUseShaders {
		return nil, fmt.Errorf("shader support has been turned off using --debug_use_shaders=false")
	}
	cacheName := vfs.Canonical("shaders", name)
	sp := shaderPath{cacheName, fmt.Sprint(params)}
	if shader, found := cache[sp]; found {
		return shader, nil
	}
	data, err := vfs.Load("shaders", name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	shaderCode, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("could not read shader %q: %v", name, err)
	}
	// Add some basic templating so we can remove branches from the shaders.
	// Not using text/template so that shader files can still be processed by gofmt.
	for name, value := range params {
		shaderCode = bytes.ReplaceAll(shaderCode, []byte("PARAMS[\""+name+"\"]"), []byte("(("+value+"))"))
	}
	shader, err := ebiten.NewShader(shaderCode)
	if err != nil {
		return nil, fmt.Errorf("could not compile shader %q: %v", name, err)
	}
	cache[sp] = shader
	return shader, nil
}
