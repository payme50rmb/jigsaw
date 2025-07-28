package logger

import (
	"fmt"
	"strings"
	"time"
)

func Default() Logger {
	return New()
}

// New create a logger with kvs
// kvs must be even
func New(kvs ...any) Logger {
	l := &logger{
		kvs: make([]kv, 0),
	}
	if len(kvs) == 0 {
		return l
	}
	if len(kvs)%2 != 0 {
		panic("kvs must be even")
	}
	for i := 0; i < len(kvs); i += 2 {
		l.kvs = append(l.kvs, kv{Key: kvs[i].(string), Value: kvs[i+1]})
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

type kv struct {
	Key   string
	Value any
}

type kvs []kv

func (k kvs) String() string {
	var b strings.Builder
	for _, kv := range k {
		b.WriteString(fmt.Sprintf("%s=%v,", kv.Key, kv.Value))
	}
	s := b.String()
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

func (k kvs) Map() map[string]interface{} {
	m := make(map[string]interface{})
	for _, _kv := range k {
		m[_kv.Key] = _kv.Value
	}
	return m
}

func (k kvs) Json() string {
	var _kvs = make([]string, 0)
	for _, _kv := range k {
		_kvs = append(_kvs, fmt.Sprintf("\"%s\":\"%v\"", _kv.Key, _kv.Value))
	}
	return fmt.Sprintf("{%s}", strings.Join(_kvs, ","))
}

type logger struct {
	kvs []kv
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
	if err == nil {
		l.Log("error", msg)
		return
	}
	l.F("error", err).Log("error", msg)
}

func (l *logger) Log(level string, msg string) {
	pmsg := make(kvs, 0)
	pmsg = append(pmsg, kv{Key: "level", Value: strings.ToUpper(level)})
	pmsg = append(pmsg, kv{Key: "msg", Value: msg})
	pmsg = append(pmsg, kv{Key: "time", Value: time.Now().Format("2006-01-02 15:04:05")})
	pmsg = append(pmsg, l.kvs...)
	fmt.Println(pmsg.Json())
}

func (l *logger) F(key string, value any) Logger {
	_l := &logger{kvs: make([]kv, 0)}
	_l.kvs = append(_l.kvs, l.kvs...)
	_l.kvs = append(_l.kvs, kv{Key: key, Value: value})
	return _l
}
