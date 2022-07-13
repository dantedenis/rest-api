package logger

import (
	"log"
	"io"
)

type Logger struct {
	ErrorMsg *log.Logger
	InfoMsg *log.Logger
}

func NewLogger(info, error io.Writer) *Logger {
	return &Logger{
		ErrorMsg: log.New(error, "\033[31m[ERROR]\t\033[0m", log.Ldate|log.Ltime|log.Lshortfile),
		InfoMsg: log.New(info, "\033[32m[INFO]\t\033[0m", log.Ldate|log.Ltime|log.Lshortfile),
	}
}