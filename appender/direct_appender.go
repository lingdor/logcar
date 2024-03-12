package appender

import (
	"io"

	"github.com/lingdor/go-logcar/entity"
	"github.com/lingdor/logcar/cfg"
	"github.com/lingdor/logcar/filter"
)

type DirectAppender struct {
	commonFilter *filter.CommonFilter
	ch           <-chan *entity.LogLine
	Writer       io.Writer
}

func (f *DirectAppender) Init(appenderConfig *cfg.AppenderConfig, ch <-chan *entity.LogLine) {

	f.commonFilter = filter.NewFilter(appenderConfig.Option)
	f.ch = ch
	go f.goConsume()
}

// Write write log content
func (f *DirectAppender) goConsume() {
	for {
		line := <-f.ch
		if f.commonFilter.IsMatch(line) {
			f.Writer.Write(line.Line)
		}
	}
}
