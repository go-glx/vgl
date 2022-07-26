package instance

import "github.com/go-glx/vgl/config"

type CreateOptions struct {
	logger             config.Logger
	appName            string
	engineName         string
	requiredExtensions []string
	debugMode          bool
}

func NewCreateOptions(
	logger config.Logger,
	appName string,
	engineName string,
	requiredExtensions []string,
	debugMode bool,
) CreateOptions {
	return CreateOptions{
		logger:             logger,
		appName:            appName,
		engineName:         engineName,
		requiredExtensions: requiredExtensions,
		debugMode:          debugMode,
	}
}
