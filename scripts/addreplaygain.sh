#!/bin/sh

find .. -path '*/music/*.ogg' | while read -r file; do
	gain=$(loudgain --maxtpl=0 --pregain=-5 --output-new "$file" | tail -n 1 | cut -f 9)
	gain=$(echo "e(l(10) * ${gain% dB} / 20)" | bc -l)
	echo "$file -> $gain"
	data=$(cat "$file.json" || echo "{}")
	echo "$data" | jq ". + {replay_gain: $gain}" > "$file.json"
done
