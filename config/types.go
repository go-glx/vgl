package config

// todo: check:
// prefix "vk:"
// suffix "\n"
// on all logs calls

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Notice(msg string)
	Error(msg string)
}
