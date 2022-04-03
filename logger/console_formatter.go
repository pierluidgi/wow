package logger

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const timeFormat = "2006-01-02 15:04:05"

type ConsoleFormatter struct {
	FilesPrefixLen int
}

//Format renders a single log entry
func (f *ConsoleFormatter) Format(entry *log.Entry) ([]byte, error) {
	var sb *bytes.Buffer
	if entry.Buffer != nil {
		sb = entry.Buffer
	} else {
		sb = &bytes.Buffer{}
	}

	sb.WriteString(entry.Time.Format(timeFormat))
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))

	if entry.HasCaller() {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%s:%d (%s)", entry.Caller.File[f.FilesPrefixLen:], entry.Caller.Line, entry.Caller.Function))
	}

	sb.WriteByte(' ')
	sb.WriteString(entry.Message)

	for k, v := range entry.Data {
		if k == "file" || k == "func" {
			continue
		}
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%s=%v", k, v))
	}

	sb.WriteByte('\n')

	return sb.Bytes(), nil
}
