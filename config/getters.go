package config

func (c *Config) InDebug() bool {
	return c.debug
}

func (c *Config) HasGPUVSync() bool {
	return c.gpu.vSync
}
