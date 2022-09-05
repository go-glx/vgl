package vlk

import (
	"reflect"
)

type (
	resource interface {
		Free()
	}

	resourceQueueCloser interface {
		EnqueueFree(fn func())
	}

	resourcesDict map[uintptr]any
)

var knownResources = resourcesDict{}

func resolver[T any](c resourceQueueCloser, factory func() *T) *T {
	// take unique ptr to golang function
	// all static function (defined in compile time), always has
	// same memory pointers. So this is unique identify for factory
	// that create this resource
	ptr := reflect.ValueOf(factory).Pointer()

	if target, exist := knownResources[ptr]; exist {
		return target.(*T)
	}

	knownResources[ptr] = factory()

	c.EnqueueFree(func() {
		if res, ok := knownResources[ptr].(resource); ok {
			res.Free()
		}

		delete(knownResources, ptr)
	})

	return knownResources[ptr].(*T)
}

func static[T any](c *Container, factory func() *T) *T {
	return resolver(c.closer, factory)
}

func dynamic[T any](c *Container, factory func() *T) *T {
	return resolver(c.rebuilder, factory)
}
