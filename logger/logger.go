package logger

import (
	"log"
	"os"
	"fmt"
	"sync"
)

type Logger struct {
	logMutex sync.Mutex
	Log      *log.Logger
	prefix   string
}

func (l *Logger) Init() {
	l.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.prefix = "\033[0;30m"
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintf(format, v...))
	l.logMutex.Unlock()
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprint(v...))
	l.logMutex.Unlock()
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintln(v...))
	l.logMutex.Unlock()
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprint(v...))
	os.Exit(1)
	l.logMutex.Unlock()
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
	l.logMutex.Unlock()
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.logMutex.Lock()
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
	l.logMutex.Unlock()
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.logMutex.Lock()
	l.prefix = prefix
	l.logMutex.Unlock()
}
