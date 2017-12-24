package logger

import (
	"log"
	"os"
)

func Logger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}