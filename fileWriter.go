package mlog

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type LogFile struct {
	filePath string
	logType  LogType
	fileName string
	file     *os.File
	maxSize  int64
}
const Suffix = ".log"
func MakeLogWriter(path string, name string, maxSize int64) (logFile *LogFile, err error) {
	logFile = &LogFile{
		fileName: name,
		filePath: path,
		maxSize:  maxSize,
	}
	err = logFile.initFile()
	if err != nil {
		return
	}
	return
}
func (this *LogFile) checkSize(size int64) bool {
	fileInfo, err := this.file.Stat()
	if err != nil {
		fmt.Println("Get File Info Failed")
		return false
	}
	return fileInfo.Size() > size
}
func (this *LogFile) newLogFile() (err error) {
	logPath := path.Join(this.filePath, this.fileName)
	file, err := os.OpenFile(logPath+Suffix, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Create New Log File Failed")
		return
	}
	this.file = file
	return
}
func (this *LogFile) initFile() (err error) {
	logPath := path.Join(this.filePath, this.fileName)
	file, err := os.OpenFile(logPath+Suffix, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open mlog file failed, err: %v\n", err)
		return
	}
	this.file = file
	return
}
func (this *LogFile) Close() (err error) {
	this.file.Close()
	if err != nil {
		fmt.Println(err)
	}
	return
}


func (this *LogFile) backLogFile() (err error) {
	newFileName := fmt.Sprintf("%s_%s%s", this.fileName, time.Now().Format("20060102150405000"),Suffix)
	err = os.Rename(path.Join(this.filePath, this.fileName+Suffix), path.Join(this.filePath, newFileName))
	if err != nil {
		log.Fatalln(err)
		fmt.Println("Back Log File Failed")
	}
	return
}
func (this *LogFile) write(logType LogType, msg string) {
	logInfo, err := getInfo(4)
	if err != nil {
		return
	}
	fmt.Println(this.checkSize(this.maxSize))
	if this.checkSize(this.maxSize) {
		this.Close()
		this.backLogFile()
		this.newLogFile()
	}
	fmt.Fprintf(this.file, "[%s] %s %d %s \n", logTypeToString(logType), logInfo.File, logInfo.Line, msg)
}
