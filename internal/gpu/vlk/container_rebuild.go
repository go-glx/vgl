package vlk

type rebuilder struct {
	queue []func()
}

func newRebuilder() *rebuilder {
	return &rebuilder{
		queue: make([]func(), 0, 32),
	}
}

func (c *rebuilder) enqueue(fn func()) {
	c.queue = append(c.queue, fn)
}

func (c *rebuilder) free() {
	for i := len(c.queue) - 1; i >= 0; i-- {
		c.queue[i]()
	}
}
