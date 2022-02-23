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

# Simple script to run a regression test demo.
# Hint: run this under Xvfb, Xdummy or similar.
# TODO: figure out how to get efficient 3D graphics in this (e.g. via Virgil?).

if [ $# -lt 4 ]; then
	echo >&2 "Usage: $0 tag required_regex 'binary with flags...' demo1.dem demo2.dem ..."
	exit 1
fi

tag=$1; shift
required_regex=$1; shift
binary=$1; shift

minor_regression=false
run_broken=false
run_finished=false

for demo in "$@"; do
	echo >&2 "Running $demo..."
	t0=$(date +%s)
	if ! $binary \
		-audio=false \
		-batch \
		-debug_check_tile_window_size \
		-debug_check_tnih_signs \
		-debug_log_file="$demo.$tag.log" \
		-demo_play="$demo" \
		-demo_record="$demo.$tag.actual.dem" \
		-demo_play_regression_prefix="$demo.$tag" \
		-demo_timedemo \
		-draw_blurs=false \
		-draw_outside=false \
		-draw_visibility_mask=false \
		-expand_using_vertices_accurately=false \
		-fps_divisor=4 \
		-fullscreen=false \
		-profiling=1m \
		-runnable_when_unfocused \
		-screen_filter=nearest \
		-show_fps \
		-show_time \
		-vsync=false \
		-window_scale_factor=1 \
		; then
		if grep -q 'regression test failed from' "$demo.$tag.log"; then
			if grep -q 'REGRESSION: difference in final save state' "$demo.$tag.log"; then
				echo "$demo had a regression that impacted save states; see log and screenshots. Probably reject?"
				run_broken=true  # Continue ahead anyway, as it is likely helpful to learn about ALL serious regressions.
			else
				echo "$demo had a regression that did not impact save states; see log and screenshots. Maybe accept?"
				minor_regression=true  # Continue ahead anyway, as this may be salvageable.
			fi
		else
			# Other cause of death.
			echo "$demo had a fatal error; see log."
			exit 4
		fi
	fi
	t1=$(date +%s)
	dt=$((t1 - t0))
	frames=$(wc -l < "$demo")
	frames=$((frames - 1))  # Deduct the "FinalSaveGame" pseudo-frame at the end.
	tps=$((frames / dt))
	echo "$demo finished after $dt seconds ($tps tps)."
	if grep -q "$required_regex" "$demo.$tag.log"; then
		run_finished=true
	fi
done

if $run_broken; then
	echo "The run is no longer complete; see above. Some manual fixes required."
	exit 3
fi

if ! $run_finished; then
	echo "The run did not end the intended way; the logs should have contained /$required_regex/. Some manual fixes required."
	exit 2
fi

if $minor_regression; then
	echo "Minor regression but run still succeeds; see above. Demos can be automatically fixed by a play+record cycle if the deltas are OK."
	exit 1
fi

echo "Run succeeded."
exit 0
