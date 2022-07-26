package config

import "log"

type defaultLogger struct {
}

func (d *defaultLogger) Debug(msg string) {
	d.log("Debug", msg)
}

func (d *defaultLogger) Info(msg string) {
	d.log("Info", msg)
}

func (d *defaultLogger) Notice(msg string) {
	d.log("Notice", msg)
}

func (d *defaultLogger) Error(msg string) {
	d.log("Error", msg)
}

func (d *defaultLogger) log(level string, msg string) {
	log.Printf("[%s] vk: %s\n", level, msg)
}
