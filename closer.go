package vgl

type Closer struct {
	queue     []func()
	backQueue []func()
}

func newCloser() *Closer {
	return &Closer{
		queue:     make([]func(), 0, 32),
		backQueue: make([]func(), 0, 4),
	}
}

func (c *Closer) EnqueueBackFree(fn func()) {
	c.backQueue = append(c.backQueue, fn)
}

func (c *Closer) EnqueueFree(fn func()) {
	c.queue = append(c.queue, fn)
}

func (c *Closer) close() {
	for i := 0; i < len(c.backQueue); i++ {
		c.backQueue[i]()
	}

	for i := len(c.queue) - 1; i >= 0; i-- {
		c.queue[i]()
	}
}
