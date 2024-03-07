package cfg

import (
	"github.com/ghodss/yaml"
	"io"
	"os"
)

var AppSet struct {
	Logcar struct {
		Listen struct {
			Port int
			Ip   string
		}
		Appenders []AppenderConfig
	}
}

type AppenderConfig struct {
	Type   string
	Filter struct {
		Levels string
		Regex string
	}
	Option map[string]any
}

func LoadConfigFile(fname string) error {
	if file, err := os.Open(fname); err == nil {
		defer file.Close()
		var bs []byte
		if bs, err = io.ReadAll(file); err == nil {
			return yaml.Unmarshal(bs, &AppSet)
		}
		return err
	}
	return nil
}

//18035905955
