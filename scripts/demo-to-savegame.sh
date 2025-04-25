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

# Script to extract savegames from demo files.

mode=$1
demo=$2

case "$mode" in
	start)
		head -n 1 < "$demo" | json_xs -e '$_ = $_->{SaveGame}'
		;;
	end)
		tail -n 1 < "$demo" | json_xs -e '$_ = $_->{FinalSaveGame}'
		;;
	diff)
		a=$(mktemp)
		b=$(mktemp)
		"$0" start "$demo" > "$a"
		"$0" end   "$demo" > "$b"
		diff -u "$a" "$b"
		status=$?
		rm -f "$a" "$b"
		exit "$status"
		;;
	diff_start)
		a=$(mktemp)
		b=$(mktemp)
		"$0" start "$demo" > "$a"
		"$0" start "$3" > "$b"
		diff -u "$a" "$b"
		status=$?
		rm -f "$a" "$b"
		exit "$status"
		;;
	diff_end)
		a=$(mktemp)
		b=$(mktemp)
		"$0" end "$demo" > "$a"
		"$0" end "$3" > "$b"
		diff -u "$a" "$b"
		status=$?
		rm -f "$a" "$b"
		exit "$status"
		;;
	diff_seq)
		a=$(mktemp)
		b=$(mktemp)
		status=0
		while [ $# -ge 3 ]; do
			"$0" end   "$2" > "$a"
			"$0" start "$3" > "$b"
			echo "--- $2.end"
			echo "+++ $3.start"
			if ! diff -u "$a" "$b"; then
				status=1
			fi
			status=$?
			shift
		done
		rm -f "$a" "$b"
		exit "$status"
		;;
	replace_start)
		newsave=$(cat)
		while read -r l; do
			case "$l" in
				*\"SaveGame\":*)
					echo "$l" | NEWSAVE=$newsave json_xs -t json -e '$_->{SaveGame} = decode_json $ENV{NEWSAVE} if exists $_->{SaveGame}'
					echo
					;;
				*)
					echo "$l"
					;;
			esac
		done < "$demo"
		;;
	*)
		echo >&2 "Usage: $0 {start|end|diff|diff_start|diff_end} filename.dem"
		echo >&2 "       $0 {replace_start} filename.dem < savegame.sav > newfilename.dem"
		echo >&2 "       $0 {diff_seq} filename.dem"
		exit 1
		;;
esac
