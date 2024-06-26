package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
var Logger *zap.Logger
func InitializeLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	//logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(os.Stdout)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
	zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
