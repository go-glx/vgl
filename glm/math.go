package glm

import "math"

const convDeg2Rad = math.Pi / circleHalfDeg
const convRad2Deg = circleHalfDeg / math.Pi
const circleHalfDeg = 180.0

func Deg2rad(deg float32) float32 {
	return deg * convDeg2Rad
}

func Rad2deg(rad float32) float32 {
	return rad * convRad2Deg
}
