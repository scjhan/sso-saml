package log

import (
	"fmt"
	"runtime"
	"time"
)

const (
	debug   = "D"
	errors  = "E"
	warning = "W"
	info    = "I"
)

func log(t string, logs string) {
	var pre string

	now := time.Now()
	if _, file, line, ok := runtime.Caller(2); ok {
		pre = fmt.Sprintf("%s [%s] [%s:%d]",
			now.Format("2006-01-02 15:04:05.000"), t, file, line)
	}

	fmt.Println(pre, logs)
}

func Debug(logs string) {
	log(debug, logs)
}

func Error(logs string) {
	log(errors, logs)
}

func Warning(logs string) {
	log(warning, logs)
}

func Info(logs string) {
	log(info, logs)
}
