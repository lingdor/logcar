package filter

import (
	"regexp"
	"strings"

	"github.com/lingdor/go-logcar/entity"
)

type CommonFilter struct {
	Levels    int
	Regex     *regexp.Regexp
	lastMatch bool
}

func NewFilter(options map[string]any) *CommonFilter {
	var ret = &CommonFilter{
		Levels: levelVals[entity.LogLevelAll],
		Regex:  nil,
	}
	if options["level"] != nil {
		ret.Levels = LevelNamesToCode(options["levels"].(string))
	}
	return ret
}

func LevelNamesToCode(names string) int {
	levels := 0
	for _, levelName := range strings.Split(names, ",") {
		switch strings.ToLower(strings.TrimSpace(levelName)) {
		case "trace":
			levels |= levelVals[entity.LogLevelTrace]
		case "debug":
			levels |= levelVals[entity.LogLevelDebug]
		case "info":
			levels |= levelVals[entity.LogLevelInfo]
		case "warn":
			levels |= levelVals[entity.LogLevelWarn]
		case "error":
			levels |= levelVals[entity.LogLevelError]
		case "fatal":
			levels |= levelVals[entity.LogLevelFatal]
		case "off":
			return entity.LogLevelOff
		case "all":
			return entity.LogLevelAll
		}
	}
	return levels
}

var levelVals = map[rune]int{
	entity.LogLevelOff:   0,
	entity.LogLevelTrace: 1,
	entity.LogLevelDebug: 2,
	entity.LogLevelInfo:  4,
	entity.LogLevelWarn:  8,
	entity.LogLevelError: 16,
	entity.LogLevelFatal: 32,
	entity.LogLevelAll:   63,
}

func (c *CommonFilter) IsMatch(line *entity.LogLine) bool {

	if line.Level == ' ' {
		return c.lastMatch
	}
	levelCode := levelVals[line.Level]
	if c.Levels&levelCode == levelCode {
		if c.Regex == nil {
			return true
		}
		return c.Regex.Match(line.Line)
	}
	return false

}
