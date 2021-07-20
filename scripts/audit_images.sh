#!/bin/sh

find . -name \*.png | sort | while read -r file; do
	# Exceptions.
	case "$file" in
		# Editing only.
		*/editorimgs/*) continue ;;
		*/src/*) continue ;;
		./assets/sprites/warpzone_*.png) continue ;;
		# Intentionally violating.
		./assets/sprites/clock_*.png) continue ;;
		./assets/sprites/gradient_*.png) continue ;;
		./assets/sprites/editorimgs/gradient_*.png) continue ;;
	esac
	set -- \
		"$file" -depth 8 +dither \
		-write MPR:orig \
		-channel RGB -remap scripts/cga_palette.pnm +channel \
		MPR:orig -alpha set -compose copy-opacity -composite \
		-channel A -threshold 25% +channel
	f=$(
		convert \
			\( "$file" -depth 8 +dither \) \
			\( "$@" \) \
			-channel RGBA \
			-metric RMSE -format '%[distortion]\n' -compare \
			INFO:
	)
	if [ "$f" !=  0 ]; then
		echo "convert "$@" "$file"  # $f"
	fi
done
