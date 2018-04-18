package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	path := fmt.Sprintf("%s/%s.%s.log", ".", "logtest2", time.Now().Format("2006.01.02"))
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
	logger := log.New(file, "", log.Llongfile)

	logger.Println("hello world")

}
