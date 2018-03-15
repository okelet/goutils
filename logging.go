package goutils

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/juju/loggo"
	"path/filepath"
	"fmt"
	"os"
)

const loggoStderrWriterName = "stderr"

func InitLogging(logPath string, logToStdErr bool) {

	if logPath == "" {
		logPath = os.DevNull
	}

	fileRotateWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    5,
		MaxBackups: 7,
	}

	loggo.ReplaceDefaultWriter(loggo.NewSimpleWriter(fileRotateWriter, func(entry loggo.Entry) string {
		return TzLoggoFormatter(entry)
	}))

	SetErrorOutputLogging(logToStdErr)

}


// Replace the default formatter with a new one that respects time zone
// https://github.com/juju/loggo/blob/master/formatter.go
func TzLoggoFormatter(entry loggo.Entry) string {
	ts := entry.Timestamp.Format("2006-01-02 15:04:05")
	filename := filepath.Base(entry.Filename)
	return fmt.Sprintf("%s %s %s %s:%d %s", ts, entry.Level, entry.Module, filename, entry.Line, entry.Message)
}


// Add a new writer to log to the error output
func SetErrorOutputLogging(enable bool) {
	if enable {
		loggo.RegisterWriter(loggoStderrWriterName, loggo.NewSimpleWriter(os.Stderr, func(entry loggo.Entry) string {
			return TzLoggoFormatter(entry)
		}))
	} else {
		loggo.RemoveWriter(loggoStderrWriterName)
	}
}
