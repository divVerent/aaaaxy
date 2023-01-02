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
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

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

func Load(name string, params interface{}) (*ebiten.Shader, error) {
	if !*debugUseShaders {
		return nil, errors.New("shader support has been turned off using --debug_use_shaders=false")
	}
	sp := shaderPath{name, fmt.Sprint(params)}
	if shader, found := cache[sp]; found {
		return shader, nil
	}
	data, err := vfs.Load("shaders", name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %w", err)
	}
	defer data.Close()
	shaderCode, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("could not read shader %q: %w", name, err)
	}
	tmpl, err := template.New("name").Delims("T[\"", "\"]").Parse(string(shaderCode))
	if err != nil {
		return nil, fmt.Errorf("could not parse shader template %q: %w", name, err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return nil, fmt.Errorf("could not execute shader template %q: %w", name, err)
	}
	shader, err := ebiten.NewShader(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not compile shader %q: %w", name, err)
	}
	cache[sp] = shader
	return shader, nil
}
