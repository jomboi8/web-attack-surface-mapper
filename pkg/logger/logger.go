package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileLogger *log.Logger

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return err
	}
	logFile := filepath.Join(logDir, fmt.Sprintf("apguard-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	fileLogger = log.New(f, "", log.LstdFlags|log.Lshortfile)
	return nil
}

func Info(format string, v ...any) {
	msg := fmt.Sprintf("[INFO] "+format, v...)
	log.Println(msg)
	if fileLogger != nil {
		fileLogger.Println(msg)
	}
}

func Warn(format string, v ...any) {
	msg := fmt.Sprintf("[WARN] "+format, v...)
	log.Println(msg)
	if fileLogger != nil {
		fileLogger.Println(msg)
	}
}

func Error(format string, v ...any) {
	msg := fmt.Sprintf("[ERROR] "+format, v...)
	log.Println(msg)
	if fileLogger != nil {
		fileLogger.Println(msg)
	}
}
