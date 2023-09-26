package system

import "log"

type DefaultLogger struct {
}

func (o DefaultLogger) Output(val string) {
	log.Println(val)
}

func NewDefaultLogger() DefaultLogger {
	return DefaultLogger{}
}
