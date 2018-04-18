package log

import (
	"container/list"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogService struct {
	path     string
	capacity int64

	file  *os.File
	queue *list.List
	lock  *sync.Mutex
}

var gLogService *LogService

func InitLogService(name string, path string) *LogService {
	if gLogService == nil {
		gLogService = &LogService{}
	}

	gLogService.path = fmt.Sprintf("%s/%s.%s.log", path, name, time.Now().Format("2006.01.02"))
	gLogService.capacity = 0

	gLogService.file, _ = os.OpenFile(gLogService.path, os.O_CREATE|os.O_APPEND, 0644)
	gLogService.queue = &list.List{}
	gLogService.lock = &sync.Mutex{}

	go gLogService.timer()

	return gLogService
}

func (s *LogService) write(logs string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.queue.PushBack(logs)
}

func (s *LogService) doTimer() {
	var str string

	s.lock.Lock()
	if s.queue.Len() != 0 {
		head := s.queue.Front()
		for head != nil {
			str = str + head.Value.(string) + "\n"
			head = head.Next()
		}
		s.queue = list.New()
	}
	s.lock.Unlock()

	if len(str) != 0 {
		if n, err := s.file.WriteString(str); err != nil {
			fmt.Println("write error")
		} else {
			fmt.Println(fmt.Sprintf("write to %s : %d bytes\n", s.file.Name(), n))
		}
	}
}

func (s *LogService) timer() {
	for {
		s.doTimer()
		time.Sleep(time.Second)
	}
}

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

	//gLogService.write(fmt.Sprintf("%s %s", pre, logs))
	fmt.Println(fmt.Sprintf("%s %s", pre, logs))
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
