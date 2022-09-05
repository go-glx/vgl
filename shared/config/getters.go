package config

import "github.com/go-glx/vgl/shared/vlkext"

func (c *Config) InDebug() bool {
	return c.debug
}

func (c *Config) IsMobileFriendly() bool {
	return c.gpu.mobileFriendly
}

func (c *Config) Logger() vlkext.Logger {
	return c.logger
}
