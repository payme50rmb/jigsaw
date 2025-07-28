package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// New create a logger with kvs
// kvs must be even
func New(kvs ...any) Logger {
	l := &logger{
		fields: make(map[string]interface{}),
	}
	if len(kvs) == 0 {
		return l
	}
	if len(kvs)%2 != 0 {
		panic("kvs must be even")
	}
	for i := 0; i < len(kvs); i += 2 {
		l.fields[kvs[i].(string)] = kvs[i+1]
	}
	return l
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Log(level string, msg string)
	F(key string, value interface{}) Logger
}

type logger struct {
	fields map[string]interface{}
}

func (l *logger) Debug(msg string) {
	l.Log("debug", msg)
}

func (l *logger) Info(msg string) {
	l.Log("info", msg)
}

func (l *logger) Warn(msg string) {
	l.Log("warn", msg)
}

func (l *logger) Error(msg string, err error) {
	l.F("error", err).Log("error", msg)
}

func (l *logger) Log(level string, msg string) {
	pmsg := l.fields
	pmsg["level"] = level
	pmsg["msg"] = msg
	pmsg["time"] = time.Now().Format("2006-01-02 15:04:05")
	b, _ := json.Marshal(pmsg)
	fmt.Println(string(b))
}

func (l *logger) F(key string, value interface{}) Logger {
	_l := &logger{
		fields: l.fields,
	}
	_l.fields[key] = value
	return l
}
