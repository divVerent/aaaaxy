#!/bin/sh
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Script to re-record a demo to fix integration tests.
in=$1; shift
out=$1; shift

savefile=$(mktemp "$HOME/.local/share/AAAAXY/save-XXXXXX.json")
trap 'rm -f "$savefile"' EXIT
savestate=${savefile##*/save-}
savestate=${savestate%.json}

# 0. Tell the user what to do.
scripts/demo-to-savegame.sh diff "$in" "$out"
echo 'OK?'
read -r ok

# 1. Extract the initial savegame.
scripts/demo-to-savegame.sh start "$in" > "$savefile"

# 2. Run the game.
"$@" -save_state="$savestate" -demo_record="$out"

# 3. Compare the final save state.
scripts/demo-to-savegame.sh diff_end "$in" "$out"
