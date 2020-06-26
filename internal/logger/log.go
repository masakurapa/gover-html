package logger

import (
	"fmt"
	"time"
)

var L Logger

type Logger struct {
	verbose bool
}

func New(v bool) {
	L = Logger{verbose: v}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	if !l.verbose {
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("[%s]%s\n", t, format), a...)
}
