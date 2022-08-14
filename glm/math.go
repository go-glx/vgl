package glm

func Clamp(v float32, min float32, max float32) float32 {
	if v < min {
		v = min
	}

	if v > max {
		v = max
	}
	
	return v
}
