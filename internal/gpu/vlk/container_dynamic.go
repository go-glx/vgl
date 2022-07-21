package vlk

import (
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
)

func (c *Container) commandPool() *command.Pool {
	return dynamic(c, &c.vlkCommandPool,
		func(x *command.Pool) { x.Free() },
		func() *command.Pool {
			return command.NewPool(
				c.physicalDevice(),
				c.logicalDevice(),
			)
		},
	)
}
