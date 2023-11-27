package logger

import (
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/util"

	"os"
)

type Logger struct {
	level int
}

var logger *Logger

func BuildLogger(level string) {
	var logLevel int
	switch level {
	case "Error":
		logLevel = consts.LevelError
	case "Warning":
		logLevel = consts.LevelWarning
	case "Info":
		logLevel = consts.LevelInfo
	case "Debug":
		logLevel = consts.LevelDebug
	}
	logvar := Logger{
		level: logLevel,
	}
	logger = &logvar
}

func Log() *Logger {
	if logger == nil {
		logvar := Logger{
			level: consts.LevelInfo,
		}
		logger = &logvar
	}
	return logger
}

// 日志写入函数
func (logvar *Logger) Panic(service string, handler string, m ...any) {
	if consts.LevelError > logvar.level {
		return
	}
	msg := fmt.Sprint("[Panic] "+"["+handler+"] ", fmt.Sprint(m...))
	util.CommonLog(service, msg)
	os.Exit(0)
}

func (logvar *Logger) Error(service string, handler string, m ...any) {
	if consts.LevelError > logvar.level {
		return
	}
	msg := fmt.Sprint("[Error] "+"["+handler+"] ", fmt.Sprint(m...))
	util.CommonLog(service, msg)
}

func (logvar *Logger) Warning(service string, handler string, m ...any) {
	if consts.LevelWarning > logvar.level {
		return
	}
	msg := fmt.Sprint("[Warning] "+"["+handler+"] ", fmt.Sprint(m...))
	util.CommonLog(service, msg)
}

func (logvar *Logger) Info(service string, handler string, m ...any) {
	if consts.LevelInfo > logvar.level {
		return
	}
	msg := fmt.Sprint("[Info] "+"["+handler+"] ", fmt.Sprint(m...))
	util.CommonLog(service, msg)
}

func (logvar *Logger) Debug(service string, handler string, m ...any) {
	if consts.LevelDebug > logvar.level {
		return
	}
	msg := fmt.Sprint("[Debug] "+"["+handler+"] ", fmt.Sprint(m...))
	util.CommonLog(service, msg)
}
