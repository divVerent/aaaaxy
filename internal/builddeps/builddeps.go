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

package builddeps

// Import some stuff we may need at build time.
// That way we include it in the repo's go.mod and go.sum file.
import (
	// To provide the "rsrc" tool.
	_ "github.com/akavel/rsrc/rsrc"

	// To provide the "go-licenses" tool.
	_ "github.com/google/go-licenses"
	_ "github.com/otiai10/copy"
	_ "github.com/spf13/cobra"
)
