package appender

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type sizeTrigger struct {
	size     int
	appender *FileAppender
	interval time.Duration
}

func (s *sizeTrigger) Init(appender *FileAppender, option map[string]any) {
	if fileSize, ok := option["fileSize"].(string); ok {
		size := strings.ToLower(fileSize)
		if len(size) > 2 && size[len(size)-1] == 'b' {
			size = size[0 : len(size)-1]
		}
		if len(size) > 1 {
			var err error
			switch size[len(size)-1] {
			case 'k':
				s.size, err = toSize(size, 1024)
			case 'm':
				s.size, err = toSize(size, 1024*1024)
			case 'g':
				s.size, err = toSize(size, 1024*1024*1024)
			case 't':
				s.size, err = toSize(size, 1024*1024*1024*1024)
			default:
				s.size, err = strconv.Atoi(size)
			}
			if err != nil {
				panic(fmt.Errorf("trigger [fileSize] value wrong:%s(%s)", size, err.Error()))
			}
		}
	}
	s.interval = time.Second
	if option["interval"] != nil {
		strv := fmt.Sprint(option["interval"])
		if v, err := strconv.Atoi(strv); err == nil {
			s.interval = time.Duration(int(time.Millisecond) * v)
		}
	}
	s.appender = appender

}
func (s *sizeTrigger) Start() {
	go s.Look()
}
func (s *sizeTrigger) Look() {
	if s.size == 0 {
		return
	}
	for {
		if stat, err := os.Stat(s.appender.Logpath); err == nil {
			if stat.Size() > int64(s.size) {
				s.appender.Achive()
			}
		}
		time.Sleep(s.interval)
	}

}

func toSize(exp string, multiply int) (int, error) {

	if v, err := strconv.Atoi(exp[0 : len(exp)-1]); err == nil {
		return v * multiply, nil
	} else {
		return 0, err
	}
}
