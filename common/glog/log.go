package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var (
	logger *zap.Logger
)

func Install() {
	// 配置编码器，使用彩色控制台编码器
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	// 日志颜色
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// 根据级别输出不同颜色
		enc.AppendString(fmt.Sprintf("\033[32m%s\033[0m", t.Format(time.RFC3339)))
	}
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	logPath := "logs"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "app.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// 创建一个按天切割日志文件的 Core
	fileCore := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(lumberjackLogger), zap.DebugLevel)

	// 创建一个将日志同时输出到控制台和文件的 Core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		fileCore,
	)

	// 创建一个带有日期时间和调用信息的 Logger
	logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	defer logger.Sync()
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(msg string, fields ...any) {
	logger.Info(fmt.Sprintf(msg, fields...))
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Warnf(msg string, fields ...any) {
	logger.Warn(fmt.Sprintf(msg, fields...))
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(msg string, fields ...any) {
	logger.Debug(fmt.Sprintf(msg, fields...))
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorf(msg string, fields ...any) {
	logger.Error(fmt.Sprintf(msg, fields...))
}
