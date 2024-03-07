package entity

import (
	"fmt"
	"testing"
)

func TestChar(t *testing.T) {

	for i := 0; i < 254; i++ {

		b := byte(i)
		c := rune(b)
		fmt.Println(b, string(c))

	}

}
