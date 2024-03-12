package appender

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/djherbis/times"
	"github.com/lingdor/go-logcar/entity"
	"github.com/lingdor/logcar/cfg"
	"github.com/lingdor/logcar/filter"
)

const YearChar = "%Y"
const MonthChar = "%m"
const DayChar = "%d"
const HourChar = "%H"
const MinuteChar = "%i"
const SecondChar = "%s"
const IndexChar = "%no"

//%Y-%m-%d-%H-%index

type FileAppender struct {
	Logpath      string
	Config       *cfg.AppenderConfig
	commonFilter *filter.CommonFilter
	ch           <-chan *entity.LogLine
	files        [2]*os.File
	filesI       *int32
	triggers     []FileTrigger
	achiveMutex  sync.Mutex
}

func (f *FileAppender) GetLogPath() string {
	return f.Logpath
}
func (f *FileAppender) CTime() (int64, error) {

	if stat, err := times.Stat(f.Logpath); err == nil {
		if stat.HasChangeTime() {
			return stat.ChangeTime().Unix(), nil
		} else if stat.HasBirthTime() {
			return stat.BirthTime().Unix(), nil
		} else {
			panic("does'nt get a ctime")
		}
	} else {
		return 0, err
	}

	// if stat, err := os.Stat(f.Logpath); err == nil {
	// 	if cstate, ok := stat.Sys().(*syscall.Stat_t); ok {
	// 		sec := cstate.Ctimespec.Sec
	// 		return sec, nil
	// 	} else {
	// 		panic(fmt.Errorf("get ctimespec faild"))
	// 	}
	// } else {
	// 	return 0, err
	// }
}

func (f *FileAppender) Init(appenderConfig *cfg.AppenderConfig, ch <-chan *entity.LogLine) {

	var i int32 = 0
	var err error
	f.filesI = &i
	f.commonFilter = filter.NewFilter(appenderConfig.Option)
	f.ch = ch
	if appenderConfig.Option["path"] == nil {
		panic(errors.New("no found path option of file-appender"))
	}
	if path, ok := appenderConfig.Option["path"].(string); ok {
		f.Logpath = path
	} else {
		panic(errors.New("path of file-appender must be string"))
	}
	if f.files[0], err = os.OpenFile(f.Logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {

		if triggers, ok := appenderConfig.Option["triggers"].([]any); ok {
			f.triggers = make([]FileTrigger, 0, len(triggers))
			for _, triggerItem := range triggers {
				if triggerMap, ok := triggerItem.(map[string]any); ok {
					var trigger FileTrigger
					if triggerMap["fileSize"] != nil {
						trigger = &sizeTrigger{}
					} else if triggerMap["period"] != nil {
						trigger = &TimeTrigger{}
					} else {
						continue
					}
					trigger.Init(f, triggerMap)
					trigger.Start()
					f.triggers = append(f.triggers, trigger)
				}
			}
		} else {
			panic(errors.New("triggers of file-appender must be array"))
		}

		go f.goConsume()
	} else {
		panic(fmt.Errorf("open error log file faild:%s", f.Logpath))
	}
}

func (f *FileAppender) Achive() {
	f.achiveMutex.Lock()
	defer f.achiveMutex.Unlock()

	if archiveOpt, ok := f.Config.Option["archive"].(map[string]any); ok {
		if tofile, ok := archiveOpt["to-file"].(string); ok {
			if ctime, err := f.CTime(); err == nil {
				ftime := time.Unix(ctime, 0)

				tofile = strings.ReplaceAll(tofile, YearChar, strconv.Itoa(ftime.Year()))
				tofile = strings.ReplaceAll(tofile, MonthChar, strconv.Itoa(int(ftime.Month())))
				tofile = strings.ReplaceAll(tofile, DayChar, strconv.Itoa(ftime.Day()))
				tofile = strings.ReplaceAll(tofile, HourChar, strconv.Itoa(ftime.Hour()))
				tofile = strings.ReplaceAll(tofile, MinuteChar, strconv.Itoa(ftime.Minute()))
				tofile = strings.ReplaceAll(tofile, SecondChar, strconv.Itoa(ftime.Second()))

				if strings.Index(tofile, IndexChar) > -1 {
					for i := 1; i < 10000000; i++ {
						ipath := strings.ReplaceAll(tofile, IndexChar, strconv.Itoa(i))
						if _, ferr := os.Stat(ipath); ferr != nil && os.IsNotExist(ferr) {
							tofile = strings.ReplaceAll(tofile, IndexChar, strconv.Itoa(i))
							break
						}
					}
				}
				os.Rename(f.Logpath, tofile)
				if newfile, err := os.OpenFile(f.Logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
					i := atomic.LoadInt32(f.filesI)
					oldf := f.files[i%2]
					f.files[(i+1)%2] = newfile
					atomic.CompareAndSwapInt32(f.filesI, i, i+1)
					oldf.Close()
				}
			}
		}
		//todo command running
	}

}

// Write write log content
func (f *FileAppender) goConsume() {
	for {
		line := <-f.ch
		if f.commonFilter.IsMatch(line) {
			var err error
			for i := 0; i < 2; i++ {
				i := atomic.LoadInt32(f.filesI) % 2
				file := f.files[i]
				if _, err = file.Write(line.Line); err == nil {
					if !line.Prefix {
						file.WriteString("\n")
					}
				} else {
					continue
				}
			}
			if err != nil {
				panic(err)
			}
		}

	}
}
