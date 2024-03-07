package appender

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	PeriodSecondly = iota
	PeriodMinutely
	PreiodHourly
	PeriodDaily
	PeriodNone
)

type TimeTrigger struct {
	appender *FileAppender
	period   int
}

func (s *TimeTrigger) Init(appender *FileAppender, option map[string]any) {
	if option["period"] != nil {
		period := strings.ToLower(fmt.Sprint(option["period"]))
		switch period {
		case "minutely", "minute":
			s.period = PeriodMinutely
			break
		case "hourly", "hour":
			s.period = PreiodHourly
			break
		case "secondly", "second":
			s.period = PeriodSecondly
			break
		case "daily", "day":
			s.period = PeriodDaily
			break
		default:
			s.period = PeriodNone
		}
	}
	s.appender = appender
}
func (c *TimeTrigger) Start() {
	go c.Look()
}
func (s *TimeTrigger) Look() {

	for {
		now := time.Now()
		var err error
		var ctime int64
		if ctime, err = s.appender.CTime(); err == nil {
			var timeLong int64
			switch s.period {
			case PeriodSecondly:
				timeLong = 1
				break
			case PeriodMinutely:
				timeLong = 60
				break
			case PreiodHourly:
				timeLong = 3600
				break
			case PeriodDaily:
				timeLong = 3600 * 24
				break
			case PeriodNone:
				return
			}
			nowTimeStamp := now.Unix()
			if nowTimeStamp-ctime < timeLong {
				sleepTime := time.Duration(timeLong-(nowTimeStamp-ctime)) * time.Second
				time.Sleep(sleepTime)
			}
			s.appender.Achive()
		} else {
			log.Printf("get file ctime faild:%s", err.Error())
			continue
		}
	}
}
