package main

import (
	"log/slog"
	"os"

	"github.com/go-logr/logr"
	"github.com/lingdor/logcar/slogcar"
)

func main() {

	logr.LogSink

	logger := slog.New(slogcar.NewJsonHandler(os.Stderr, nil))

	logger.Info("good\n123")
	logger.WithGroup("xx").Error("world", slog.Int("ff", 88))
}
