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

# Simple script to run a benchmark demo.

if [ $# -lt 2 ]; then
	echo >&2 "Usage: $0 benchmark.dem binary with flags..."
	exit 1
fi

demo=$1; shift

exec "$@" \
	-auto_adjust_quality=false \
	-batch \
	-demo_play="$demo" \
	-demo_timedemo \
	-runnable_when_unfocused \
	-vsync=false
