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
	exit 1
fi

demo=$1; shift
count=$1; shift

run() {
	tag=$1; shift
	time=$(sh scripts/measure-timedemo.sh "$demo" "$@" -load_config=false)
	if [ $? -eq 0 ]; then
		echo "$tag,$time"
	fi
}

echo "settings,start_time_sec,play_time_sec,quit_time_sec,total_time_sec"
for i in $(seq 1 "$count"); do
	run lowest "$@" -palette=none -draw_blurs=false -draw_outside=false -expand_using_vertices_accurately=false -screen_filter=nearest
	run low    "$@" -palette=none -draw_blurs=false -draw_outside=false -expand_using_vertices_accurately=true  -screen_filter=nearest
	run medium "$@" -palette=none -draw_blurs=true  -draw_outside=false -expand_using_vertices_accurately=true  -screen_filter=simple
	run high   "$@" -palette=none -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=simple
	run max    "$@" -palette=none -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=linear2xcrt
	run vga    "$@" -palette=vga  -draw_blurs=true  -draw_outside=true  -expand_using_vertices_accurately=true  -screen_filter=linear2xcrt
done
