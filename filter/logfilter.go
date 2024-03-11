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
			levels |= levelVals[entity.LogLeveLTrace]
			break
		case "debug":
			levels |= levelVals[entity.LogLeveLDebug]
			break
		case "info":
			levels |= levelVals[entity.LogLeveLInfo]
			break
		case "warn":
			levels |= levelVals[entity.LogLeveLWarn]
			break
		case "error":
			levels |= levelVals[entity.LogLeveLError]
			break
		case "fatal":
			levels |= levelVals[entity.LogLeveLFatal]
			break
		}
	}
	return levels
}

var levelVals = map[rune]int{
	entity.LogLeveLOff:   0,
	entity.LogLeveLTrace: 1,
	entity.LogLeveLDebug: 2,
	entity.LogLeveLInfo:  4,
	entity.LogLeveLWarn:  8,
	entity.LogLeveLError: 16,
	entity.LogLeveLFatal: 32,
	entity.LogLevelAll:   63,
}

func (c *CommonFilter) IsMatch(line *entity.LogLine) bool {

	if line.Prefix {
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
