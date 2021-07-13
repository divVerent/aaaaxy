#!/bin/sh

find .. -name \*.png | sort | while read -r file; do
	# Exceptions.
	case "$file" in
		# Editing only.
		*/src/*) continue ;;
		../assets/sprites/warpzone_*.png) continue ;;
		# Intentionally violating.
		../assets/sprites/clock_*.png) continue ;;
		../assets/sprites/gradient_*.png) continue ;;
		../assets/sprites/editorimgs/gradient_*.png) continue ;;
	esac
	f=$(
		convert \
			\( \
				"$file" -depth 8 -alpha off \
			\) \
			\( \
				"$file" -depth 8 -alpha off +dither \
				-channel RGB -remap cga_palette.pnm \
				-channel A -threshold 50% \
				+channel \
			\) \
			-channel RGBA \
			-metric RMSE -format '%[distortion]\n' -compare \
			INFO:
	)
	if [ "$f" !=  0 ]; then
		echo "convert \( '$file' -depth 8 -alpha off +dither -remap cga_palette.pnm \) \( '$file' -depth 8 -alpha extract -threshold 50% \) -compose CopyOpacity -composite "$file"  # $f"
	fi
done
