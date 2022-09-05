package vlkext

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Notice(msg string)
	Error(msg string)
}
