package logger

import (
	"fmt"
	"os"
)

var logLevel = "info"

var logFormat = "json"

var logOutput = "stdout"
var logFilePath = ""

var openedFile *os.File

func SetLevel(l string) {
	logLevel = l
}

func SetFormat(f string) {
	logFormat = f
}

func SetOutput(o string) {
	logOutput = o
}

func SetFilePath(p string) {
	logOutput = "file"
	logFilePath = p
}

func GetLevel() string {
	return logLevel
}

func GetFormat() string {
	return logFormat
}

func GetOutput() string {
	return logOutput
}

func printable(level string) bool {
	switch logLevel {
	case "debug":
		return true
	case "info":
		return level == "info" || level == "warn" || level == "error"
	case "warn":
		return level == "warn" || level == "error"
	case "error":
		return level == "error"
	}
	return true
}

func printTo(l string) {
	switch logOutput {
	case "stdout":
		fmt.Println(l)
		return
	case "file":
		if logFilePath == "" {
			return
		}
		if openedFile == nil {
			f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			openedFile = f
		}
		_, _ = openedFile.WriteString(l + "\n")
		return
	}
}
