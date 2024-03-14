package app

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lingdor/go-logcar/entity"
)

func readLine(reader *bufio.Reader) (*entity.LogLine, error) {
	var line []byte
	var lastprefix, prefix = false, false
	var err error
	if line, prefix, err = reader.ReadLine(); err == nil {
		if line == nil {
			return nil, nil
		}
		// fmt.Printf("%q\n", line)
		lineInf := &entity.LogLine{
			Line:   line,
			Prefix: prefix,
			Level:  ' ',
		}
		if !lastprefix {
			if len(line) > 0 && line[0] == entity.WrapChar {
				lineInf.Line = lineInf.Line[1:]
			} else if len(line) > 0 && line[0] < entity.WrapChar && line[0] > entity.LogLevelOff {
				lineInf.Level = rune(line[0])
				lineInf.Line = lineInf.Line[1:]
			} else {
				lineInf.Level = entity.LogLevelInfo
			}
		}
		// fmt.Printf("%+v,%q\n", lineInf, lineInf.Line)
		lastprefix = prefix
		newl := make([]byte, len(lineInf.Line))
		copy(newl, lineInf.Line)
		lineInf.Line = newl
		return lineInf, nil
	}
	return nil, err
}

func startStdin(chs []chan *entity.LogLine) {

	reader := bufio.NewReader(os.Stdin)
	for {
		if line, err := readLine(reader); err == nil {
			if line == nil {
				fmt.Println("done")
				return
			}

			// fmt.Printf("%+v,%q\n", line, string(line.Line))
			for _, ch := range chs {
				ch <- line
			}
		} else if err == io.EOF {
			return
		} else {
			panic(err)
		}
	}

}
