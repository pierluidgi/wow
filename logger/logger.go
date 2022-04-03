package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
)

func Init(logLevel string) {
	level, err := log.ParseLevel(logLevel)

	if err != nil {
		level = log.TraceLevel
	}

	log.SetLevel(level)
	log.SetReportCaller(true)

	filesPrefixLen := 0

	if _, f, _, ok := runtime.Caller(1); ok {
		filesPrefixLen = len(filepath.Dir(f)) + 1
	}

	log.SetFormatter(&ConsoleFormatter{
		FilesPrefixLen: filesPrefixLen,
	})

	log.SetOutput(os.Stderr)
}
