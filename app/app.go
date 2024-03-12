package app

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lingdor/go-logcar/entity"
	"github.com/lingdor/logcar/appender"
	"github.com/lingdor/logcar/cfg"
)

var argFile string

func Prepare() {

	flag.StringVar(&argFile, "f", "./logcar.yml", "config file path of logcar")

	flag.Parse()
	checkValid()
}

func checkValid() {

}

func readLine(reader *bufio.Reader) (*entity.LogLine, error) {
	var line []byte
	var prefix bool
	var err error
	if line, prefix, err = reader.ReadLine(); err == nil {
		if line == nil {
			return nil, nil
		}
		lineInf := &entity.LogLine{
			Line:   line,
			Prefix: prefix,
			//			Level:  entity.LogLeveInfo,
		}
		if !prefix {
			if len(line) > 0 && line[0] <= entity.WrapChar && line[0] >= entity.LogLeveLOff {
				lineInf.Level = rune(line[0])
				lineInf.Line = lineInf.Line[1:]
				lineInf.Prefix = (line[0] == entity.WrapChar)
			}
		}
		return lineInf, nil
	}
	return nil, err
}

func StartApp() {

	var err error

	if err = cfg.LoadConfigFile(argFile); err == nil {
		appenderCnt := len(cfg.AppSet.Logcar.Appenders)
		var chs = make([]chan *entity.LogLine, appenderCnt)
		var appenders = make([]appender.Appender, appenderCnt)
		for i := 0; i < appenderCnt; i++ {
			chs[i] = make(chan *entity.LogLine, 100)
			appenders[i] = newAppender(&cfg.AppSet.Logcar.Appenders[i], chs[i])
		}
		//read/io
		reader := bufio.NewReader(os.Stdin)
		for {
			if line, err := readLine(reader); err == nil {
				if line == nil {
					return
				}
				for _, ch := range chs {
					ch <- line
				}
			}
		}
	}
	if err != nil {
		panic(err)
	}
}

func newAppender(appenderConfig *cfg.AppenderConfig, ch <-chan *entity.LogLine) appender.Appender {

	var ret appender.Appender

	switch strings.ToLower(appenderConfig.Type) {
	case "file":
		logpath := appenderConfig.Option["path"]
		if logpath == nil || logpath.(string) == "" {
			panic(fmt.Errorf("no found file-appender path option"))
		}
		ret = &appender.FileAppender{
			Logpath: logpath.(string),
			Config:  appenderConfig,
		}
	case "stdout":
		ret = &appender.DirectAppender{
			Writer: os.Stdout,
		}
	case "stderr":
		ret = &appender.DirectAppender{
			Writer: os.Stderr,
		}
	default:
		panic(fmt.Errorf("no found appender type:%s", appenderConfig.Type))
	}
	ret.Init(appenderConfig, ch)
	return ret
}
