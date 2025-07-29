package logger

import "fmt"

func Info(msg string, kvs ...any) {
	New(kvs...).Info(msg)
}

func Error(msg string, err error, kvs ...any) {
	New(kvs...).Error(msg, err)
}

func Debug(msg string, kvs ...any) {
	New(kvs...).Debug(msg)
}

func Warn(msg string, kvs ...any) {
	New(kvs...).Warn(msg)
}

func Errorf(format string, a ...any) {
	Default().Error(fmt.Sprintf(format, a...), nil)
}

func Infof(format string, a ...any) {
	Default().Info(fmt.Sprintf(format, a...))
}

func Debugf(format string, a ...any) {
	Default().Debug(fmt.Sprintf(format, a...))
}

func Warnf(format string, a ...any) {
	Default().Warn(fmt.Sprintf(format, a...))
}
