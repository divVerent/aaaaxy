#!/bin/sh

name=$1; shift
file=$1; shift

start=$(date +%s.%N)
"$@"
status=$?
end=$(date +%s.%N)

# Subtract them but we don't always have bc.
start_s=${start%.*}
start_ns=1${start##*.}
end_s=${end%.*}
end_ns=1${end##*.}
if [ "$start_ns" -gt "$end_ns" ]; then
	delta_ns=$(( end_ns + 1000000000 - start_ns ))
	delta_s=$(( end_s - 1 - start_s ))
else
	delta_ns=$(( end_ns - start_ns ))
	delta_s=$(( end_s - start_s ))
fi
delta=$(printf %d.%09d "$delta_s" "$delta_ns")

cat >"$file" <<EOF
[
    {
        "name": "$name",
        "unit": "Seconds",
        "value": $delta
    }
]
EOF
exit "$status"
