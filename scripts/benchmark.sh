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

# Simple script to run a benchmark demo at all graphics levels.

if [ $# -lt 2 ]; then
	echo >&2 "Usage: $0 benchmark.dem count binary with flags..."
	echo >&2 'Optionally, export AAAAXY_FRAME_PROFILING=filename to gather logs.'
	exit 1
fi

demo=$1; shift
count=$1; shift

measure() {
	(
		if [ -n "$AAAAXY_FRAME_PROFILING" ]; then
			exec 3>&2
			exec 2>"$AAAAXY_FRAME_PROFILING.$tag.log"
			set -- "$@" -debug_frame_profiling
		fi
		sh scripts/measure-timedemo.sh "$demo" "$@" -load_config=false
		if [ -n "$AAAAXY_FRAME_PROFILING" ]; then
			exec 2>&3
			exec 3>&-
			perl -ne '
					/frame time: (?:([0-9]+)h)?(?:([0-9]+)m)?(?:([0-9.]+)s)?(?:([0-9.]+)ms)?(?:([0-9.]+)µs)?(?:([0-9.]+)ns)?$/
						or next;
					my $t = $1 * 3600 + $2 * 60 + $3 + $4 * 1e-3 + $5 * 1e-6 + $6 * 1e-9;
					print "$t\n";
				' \
				< "$AAAAXY_FRAME_PROFILING.$tag.log" \
				> "$AAAAXY_FRAME_PROFILING.$tag.plot"
		fi
	)
}

run() {
	tag=$1; shift
	time=$(measure "$@")
	if [ $? -eq 0 ]; then
		echo "$tag,$time"
	fi
}

echo "settings,start_time_sec,play_time_sec,quit_time_sec,total_time_sec"
for i in $(seq 1 "$count"); do
	run lowest "$@" -palette=none -draw_blurs=false -draw_outside=false -expand_using_vertices_accurately=false -screen_filter=nearest
	run low    "$@" -palette=none -draw_blurs=false -draw_outside=false -expand_using_vertices_accurately=true  -screen_filter=nearest
	run medium "$@" -palette=none -draw_blurs=true  -draw_outside=false -expand_using_vertices_accurately=true  -screen_filter=linear2x
	run high   "$@" -palette=none -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=linear2x
	run max    "$@" -palette=none -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=linear2xcrt
	run vga    "$@" -palette=vga  -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=linear2xcrt
done
