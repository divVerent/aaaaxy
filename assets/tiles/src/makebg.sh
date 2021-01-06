#!/bin/sh

colors='
	8#555555
	9#5555ff
	a#55ff55
	b#55ffff
	c#ff5555
	d#ff55ff
	e#ffff55
	f#ffffff
'

for kc1 in $colors; do
	k1=${kc1%%#*}
	c1=${kc1#?}
	for kc2 in $colors; do
		k2=${kc2%%#*}
		c2=${kc2#?}
		convert bgtransition.png \
			-fill "$c1" -opaque "#00fe00" \
			-fill "$c2" -opaque "#ff00fe" \
			../bg_"$k1$k2"_v.png
		convert bgtransition.png \
			-fill "$c1" -opaque "#00fe00" \
			-fill "$c2" -opaque "#ff00fe" \
			-rotate 270 \
			../bg_"$k1$k2"_h.png
	done
	convert -size 16x16 xc:"$c1" ../bg_"$k1".png
done
