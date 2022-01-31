#!/bin/sh
# Copyright 2021 Google LLC
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

# Simple script to play a list of demos, and reconnect them by their savegames.
# Only works if the same inputs still do the same - i.e. the demos should have
# passed the regression test other than for savegame differences that do not
# impact gameplay at all.

set -ex

save=
for demo in "$@"; do
	if [ -n "$save" ]; then
		echo "$save" | scripts/demo-to-savegame.sh replace_start "$demo" > "$demo.replaced.dem"
	else
		cp "$demo" "$demo.replaced.dem"
	fi
	./aaaaxy \
		-audio=false \
		-batch \
		-debug_enable_drawing=false \
		-demo_play="$demo.replaced.dem" \
		-demo_record="$demo.rechained.dem" \
		-demo_timedemo \
		-fullscreen=false \
		-runnable_when_unfocused \
		-window_scale_factor=1 || true
	save=$(scripts/demo-to-savegame.sh end "$demo.rechained.dem")
	[ -n "$save" ]
done
