package vkconv

func ClampUint(n, min, max uint32) uint32 {
	if n <= min {
		return min
	}

	if n >= max {
		return max
	}

	return n
}
