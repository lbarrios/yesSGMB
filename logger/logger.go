package logger

import (
	"log"
	"os"
	"fmt"
)

type Logger struct {
	Log    *log.Logger
	prefix string
}

func (l *Logger) Init() {
	l.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.prefix = "\033[0;30m"
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.Log.SetPrefix(l.prefix)
	l.Log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}
