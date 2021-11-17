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

if [ $# -lt 3 ]; then
	echo >&2 "Usage: $0 tag 'binary with flags...' demo1.dem demo2.dem ..."
	exit 1
fi

tag=$1; shift
binary=$1; shift

set -x

for demo in "$@"; do
	if ! $binary \
		-audio=false \
		-batch \
		-demo_play="$demo" \
		-demo_play_regression_prefix="$demo.$tag" \
		-demo_timedemo \
		-draw_blurs=false \
		-draw_outside=false \
		-draw_visibility_mask=false \
		-expand_using_vertices_accurately=false \
		-fullscreen=false \
		-profiling \
		-runnable_when_unfocused \
		-screen_filter=simple \
		-show_fps \
		-show_time \
		-vsync=false \
		-window_scale_factor=1 \
		>"$demo.$tag.log" \
		2>&1; then
		cat "$demo.$tag.log"
		exit 1
	fi
done
