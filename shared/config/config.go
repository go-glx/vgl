package config

import "github.com/go-glx/vgl/shared/vlkext"

type (
	Config struct {
		debug  bool
		gpu    configSwapChain
		logger vlkext.Logger
	}

	configSwapChain struct {
		mobileFriendly bool
	}

	Configure = func(*Config)
)

func NewConfig(opts ...Configure) *Config {
	cfg := &Config{
		debug: false,
		gpu: configSwapChain{
			mobileFriendly: true,
		},
		logger: &defaultLogger{},
	}

	for _, configure := range opts {
		configure(cfg)
	}

	return cfg
}

// WithDebug will print vulkan validation errors
// on stdout. Its require vulkan SDK to work
func WithDebug(enabled bool) Configure {
	return func(config *Config) {
		config.debug = enabled
	}
}

// WithMobileFriendly will use FIFO rendering
// true - vsync, good for mobile (small power consumption)
// false - low latency, high power consumption, but better latency
func WithMobileFriendly(enabled bool) Configure {
	return func(config *Config) {
		config.gpu.mobileFriendly = enabled
	}
}

// WithLogger allow to use custom logger
// for library messages. If not set, default go
// log.* package will be used for logging
func WithLogger(logger vlkext.Logger) Configure {
	return func(config *Config) {
		config.logger = logger
	}
}
