package config

func (c *Config) InDebug() bool {
	return c.debug
}

func (c *Config) IsMobileFriendly() bool {
	return c.gpu.mobileFriendly
}

func (c *Config) Logger() Logger {
	return c.logger
}
