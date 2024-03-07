package filter

import (
	"log"
	"testing"

	"github.com/lingdor/logcar/entity"
)

func TestNamesLevel(t *testing.T) {

	levels := LevelNamesToCode("trace, debug, info,warn,error, fatal") //,
	if levels != levelVals[entity.LogLevelAll] {
		t.Fail()
	}

}

func TestFF(t *testing.T) {

	log.Println("123")

}
