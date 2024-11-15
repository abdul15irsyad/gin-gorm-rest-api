package services

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (ls *LogService) Init() {
	date := time.Now().Format("2006-01-02")

	logFileName := fmt.Sprintf("logs/%s.log", date)

	logFile := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}

	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
}

func (ls *LogService) WithFields(fields logrus.Fields, message string) {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.WithFields(fields).Info(message)
}

func (ls *LogService) Error(message string) {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.Error(message)
}
