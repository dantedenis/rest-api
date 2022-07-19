package logger

import (
	"io"
	"log"
)

type Logg interface {
	PrintError(...string)
	PrintfInfo(string, ...any)
}

type logger struct {
	ErrorMsg *log.Logger
	InfoMsg  *log.Logger
}

func NewLogger(info, error io.Writer) *logger {
	return &logger{
		ErrorMsg: log.New(error, "\033[31m[ERROR]\t\033[0m", log.Ldate|log.Ltime|log.Lshortfile),
		InfoMsg:  log.New(info, "\033[32m[INFO]\t\033[0m", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) PrintError(msg ...string) {
	l.ErrorMsg.Println(msg)
}

func (l *logger) PrintfInfo(msg string, a ...any) {
	l.InfoMsg.Printf(msg, a)
}
