package vlk

func static[T any](c *Container, target **T, down func(*T), up func() *T) *T {
	// already created
	if *target != nil {
		return *target
	}

	// up
	*target = up()

	// down
	c.closer.EnqueueFree(func() {
		down(*target)
		*target = nil
	})

	// return created
	return *target
}

// dynamic is same as static, but down function enqueued to rebuilder
// instead of global closer. Rebuilder can be called many times
// in engine run, for example on window resize event. This will
// break and free all dynamic resources, like graphics pipelines
// and next lazy call should rebuild this from scratch
func dynamic[T any](c *Container, target **T, down func(*T), up func() *T) *T {
	// already created
	if *target != nil {
		return *target
	}

	// up
	*target = up()

	// down
	c.rebuilder.enqueue(func() {
		down(*target)
		*target = nil
	})

	// return created
	return *target
}
