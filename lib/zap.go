package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LibLogger struct {
	Logger *zap.Logger
}

func NewLogger() *LibLogger {
	// Configure file logging
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02")),
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel)

	// Configure console logging
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	// Combine cores for multi-output logging
	core := zapcore.NewTee(fileCore, consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	return &LibLogger{Logger: logger}
}
