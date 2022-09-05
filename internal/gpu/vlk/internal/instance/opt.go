package instance

import (
	"unsafe"

	"github.com/go-glx/vgl/shared/vlkext"
)

type CreateOptions struct {
	logger             vlkext.Logger
	procAddr           unsafe.Pointer
	appName            string
	engineName         string
	requiredExtensions []string
	debugMode          bool
}

func NewCreateOptions(
	logger vlkext.Logger,
	procAddr unsafe.Pointer,
	appName string,
	engineName string,
	requiredExtensions []string,
	debugMode bool,
) CreateOptions {
	return CreateOptions{
		logger:             logger,
		procAddr:           procAddr,
		appName:            appName,
		engineName:         engineName,
		requiredExtensions: requiredExtensions,
		debugMode:          debugMode,
	}
}
