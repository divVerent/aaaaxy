// Copyright 2021 Google LLC
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

//go:build !wasm
// +build !wasm

package alert

import (
	"github.com/ncruces/zenity"
	"os/exec"
)

func Show(msg string) {
	err := zenity.Error(msg, zenity.Title("AAAAXY - Error"))
	if err != nil {
		err = exec.Command("gxmessage", "-center", "-title", "AAAAXY - Error", " "+msg).Run()
	}
	if err != nil {
		err = exec.Command("xmessage", "-center", "-title", "AAAAXY - Error", " "+msg).Run()
	}
	// No further fallback - we already wrote the fatal error to stderr anyway.
}
