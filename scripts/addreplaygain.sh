#!/bin/sh

find assets third_party -path '*/music/*.ogg' | while read -r file; do
	silence=$(sox "$file" -t dat - | awk 'BEGIN { n = 0; } /^\s*[0-9]/ { if ($2 != 0 || $3 != 0) { nextfile; } ++n; } END { print(n); }')
	gain=$(loudgain --maxtpl=0 --pregain=-5 --output-new "$file" | tail -n 1 | cut -f 9)
	gain=$(echo "e(l(10) * ${gain% dB} / 20)" | bc -l)
	echo "$file -> $gain, $silence"
	data=$(cat "$file.json" || echo "{}")
	echo "$data" | jq ". + {play_start: $silence, replay_gain: $gain}" > "$file.json"
done
