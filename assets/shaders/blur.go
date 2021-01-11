// A simple shader to perform blurs.
package main

var Size float
var Step vec2
var Scale float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	acc := imageSrc0UnsafeAt(texCoord)
	for y := 1.0; y <= 8.0; y += 1.0 {
		if y <= Size {
			d := y * Step
			acc += imageSrc0At(texCoord - d)
			acc += imageSrc0At(texCoord + d)
		}
	}
	return acc * Scale
}
