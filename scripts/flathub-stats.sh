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

# Viewer for stats from Flathub.
# Shows number of people updating to each release as a measure for total users.

rm -f stats.csv

days_back=90
days_nocache=2
id=io.github.divverent.aaaaxy
cache=$HOME/.cache/AAAAXY/flathub-stats

mkdir -p "$cache"

{
	git tag | grep '^v.*\..*\..*$' | while read -r tag; do
		date=$(git rev-list -1 --date='format:%Y-%m-%d' --pretty='format:date %ad' "$tag" | grep '^date ' | cut -d ' ' -f 2)
		echo "\"$date\",\"R\",\"$tag\""
	done
	for i in $(seq 0 $days_back); do
		date=$(date +%Y-%m-%d -d"$i days ago")
		urldate=$(echo "$date" | tr - /)
		if [ -s "$cache/$date.json" ]; then
			cat "$cache/$date.json"
		elif [ $i -lt $days_nocache ]; then
			curl "https://flathub.org/stats/$urldate.json"
		else
			trap 'rm -f "$cache/$date.json"' EXIT
			curl "https://flathub.org/stats/$urldate.json" | tee "$cache/$date.json"
			trap - EXIT
		fi | DATE=$date jq --raw-output '.refs["'"$id"'"] | [(0, 1) as $i | [(keys | .[]) as $key | .[$key][$i]] | add] | [$ENV["DATE"], "S", .[0] - .[1], .[1]] | @csv'
	done
} | sort | tr -d '"' | {
	started=false
	release=
	release_date=
	release_updates=0
	release_started=false
	echo '1 set timefmt "%Y-%m-%d"'
	echo '1 set xdata time'
	echo '1 set key left top'
	echo '1 plot "-" using 1:2 with lines title "new installs", "-" using 1:2 with lines title "update installs", "-" using 1:2 with lines title "update installs to release", "-" using 1:2:3 with labels offset char 0, -1 notitle'
	echo '3 e'
	echo '5 e'
	echo '7 e'
	echo '9 e'
	finish_release() {
		if $release_started && [ $release_updates -gt 0 ]; then
			echo "6 $release_date $release_updates"
			echo "6 $release_enddate $release_updates"
			releade_middate=$(date +'%Y-%m-%d' -d@$((($(date +%s -d"$release_date") + $(date +%s -d"$release_enddate")) / 2)))
			release_str=$(echo "$release" | sed -e 's,+,\\\\n+,g')
			echo "8 $releade_middate $release_updates \"$release_str\""
		fi
	}
	while IFS=',' read date type x y; do
		release_enddate=$date
		case "$type" in
			R)
				finish_release
				release=$x
				release_date=$date
				release_updates=0
				release_started=$started
				;;
			S)
				started=true
				installs=$x
				updates=$y
				release_updates=$((release_updates + updates))
				echo "2 $date $installs"
				echo "4 $date $updates"
				;;
		esac
	done
	finish_release
} | sort -s -k '1.1,1.1' | cut -c 3- | gnuplot -persist
