package appender

import (
	"github.com/lingdor/go-logcar/entity"
	"github.com/lingdor/logcar/cfg"
)

type Appender interface {
	Init(appenderConfig *cfg.AppenderConfig, ch <-chan *entity.LogLine)
}

type FileTrigger interface {
	Init(appender *FileAppender, option map[string]any)
	Start()
}

type FileArchive interface {
	Init(appender *FileAppender, option map[string]any)
	execute(appender *FileAppender) error
}
