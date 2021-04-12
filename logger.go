package mlog

import (
	"errors"
	"fmt"
	"path"
	"runtime"
)


//const DateFormat = "2006-01-02 15:04:05"
type LogType uint8

type Logger struct {
	level LogType
	group map[LogType]*LogFile
}
type Config struct {
	LogPath     string
	fileMaxSize int64
	Level       LogType
	TypeMapFile map[LogType]string
}

const (
	INFO LogType = 1 << iota
	DEBUG
	ERROR
	FATAL
)

func logTypeToString(logType LogType) string {
	switch logType {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type RuntimeInfo struct {
	File     string
	FuncName string
	Line     int
}

func getInfo(n int) (info RuntimeInfo, err error) {
	pc, file, line, ok := runtime.Caller(n)
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file)
	if !ok {
		err = errors.New("")
	}
	info = RuntimeInfo{
		FuncName: funcName,
		File:     fileName,
		Line:     line,
	}
	return
}



func NewLogger(config Config) (l *Logger, err error) {
	l = &Logger{
		level: config.Level,
		group: map[LogType]*LogFile{},
	}
	logFileMap := map[string]*LogFile{}
	for _, item := range []LogType{INFO, DEBUG, ERROR, FATAL} {
		if item&config.Level == item {
			logFile, ok := logFileMap[config.TypeMapFile[item]]
			fileName := config.TypeMapFile[item]
			if !ok {
				l.group[item], err = MakeLogWriter(config.LogPath, fileName, config.fileMaxSize)
				if err != nil {
					return
				}
				logFileMap[fileName] = l.group[item]
			} else {
				l.group[item] = logFile
				logFileMap[fileName] = logFile
			}
		}
	}
	return
}
func (logger *Logger) ignore(logType LogType) bool {
	if logType&logger.level == logType {
		return false
	}
	return true

}
func (logger *Logger) log(logType LogType, msg string) {
	if logger.ignore(logType) {
		return
	}
	logger.group[logType].write(logType, msg)
}
func Debug(format string, args ...interface{}) {
	Log.log(DEBUG, fmt.Sprintf(format, args...))
}
func Info(format string, args ...interface{}) {
	Log.log(INFO, fmt.Sprintf(format, args...))
}
func Error(format string, args ...interface{}) {
	Log.log(ERROR, fmt.Sprintf(format, args...))
}
func Fatal(format string, args ...interface{}) {
	Log.log(FATAL, fmt.Sprintf(format, args...))
}
