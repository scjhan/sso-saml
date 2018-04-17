package main

import (
	"fmt"
	"runtime"
	"time"
)

func Log() {
	now := time.Now()
	if _, file, line, ok := runtime.Caller(1); ok {
		fmt.Println(fmt.Sprintf("%s [D] [%s:%d]",
			now.Format("2006-01-02 15:04:05.000"), file, line))
	}
}

func main() {
	fmt.Println(time.Now().String())
	now := time.Now()

	y := now.Year()
	m := now.Month()
	d := now.Day()
	h := now.Hour()
	mi := now.Minute()
	ns := now.Nanosecond()

	now.Month()

	fmt.Printf("%d/%d/%d %d:%d:%.3f\n", y, m, d, h, mi, 1.0*ns/1000000)
	fmt.Println(now.Format("2006-01-02 15:04:05"))

	Log()
}
