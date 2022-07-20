package instance

type CreateOptions struct {
	appName            string
	engineName         string
	requiredExtensions []string
	debugMode          bool
}

func NewCreateOptions(
	appName string,
	engineName string,
	requiredExtensions []string,
	debugMode bool,
) CreateOptions {
	return CreateOptions{
		appName:            appName,
		engineName:         engineName,
		requiredExtensions: requiredExtensions,
		debugMode:          debugMode,
	}
}
