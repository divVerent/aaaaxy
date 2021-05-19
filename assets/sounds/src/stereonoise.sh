#!/bin/sh

sox \
	-r 44.1k -b 16 -c 2 -n \
	../stereonoise0.wav \
	synth 10.0 whitenoise \
	sinc -a 40 -1000 \
	gain -n -1

sox ../stereonoise0.wav \
	../sterenoise1.wav \
	trim 0.0 5.0 \
	fade p 5.0

sox ../stereonoise0.wav \
	../sterenoise2.wav \
	trim 5.0 \
	fade p 0.0 5.0 5.0

sox -m ../sterenoise1.wav ../sterenoise2.wav \
	../stereonoise.wav \

oggenc -q3 -o ../stereonoise.ogg ../stereonoise.wav
