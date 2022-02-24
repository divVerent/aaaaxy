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

# Simple script to run a benchmark demo and measure its times.

if [ $# -lt 2 ]; then
	echo >&2 "Usage: $0 benchmark.dem binary with flags..."
	exit 1
fi

demo=$1; shift
t_starting=$(date +%s.%N)
out=$(scripts/run-timedemo.sh "$@" 2>&1 | tee /dev/stderr)
status=$?
t_started=$(date +%s.%N -d"$(echo "$out" | awk '/ \[INFO\] game started$/     { print $1, $2 }')")
t_exiting=$(date +%s.%N -d"$(echo "$out" | awk '/ \[INFO\] exiting normally$/ { print $1, $2 }')")
t_exited=$(date +%s.%N)
{
	echo "$t_started * 1 - $t_starting" | bc
	echo ,
	echo "$t_exiting * 1 - $t_started" | bc
	echo ,
	echo "$t_exited * 1 - $t_exiting" | bc
	echo ,
	echo "$t_exited * 1 - $t_starting" | bc
} | tr -d '\n'
exit $status
