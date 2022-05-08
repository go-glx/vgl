package closer

import (
	"log"
)

type (
	CloseFn func() error
	FreeFn  func()
	Closer  struct {
		queue []CloseFn
	}
)

func NewCloser() *Closer {
	return &Closer{
		queue: make([]CloseFn, 0),
	}
}

func (c *Closer) EnqueueClose(fn CloseFn) {
	c.queue = append(c.queue, fn)
}

func (c *Closer) EnqueueFree(fn FreeFn) {
	c.EnqueueClose(func() error {
		fn()

		return nil
	})
}

func (c *Closer) Close() error {
	for i := len(c.queue) - 1; i >= 0; i-- {
		closeFn := c.queue[i]
		err := closeFn()
		if err != nil {
			log.Printf("close err #%d: %v\n", i, err)
		}
	}

	return nil
}
