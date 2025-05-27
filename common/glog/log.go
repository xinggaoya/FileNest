package glog

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

// ColorTimeEncoder 控制台时间字段绿色高亮
func ColorTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("\033[32m" + t.Format("2006-01-02 15:04:05.000") + "\033[0m")
}

// GrayCallerEncoder 控制台文件:行号灰色显示
func GrayCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if caller.Defined {
		enc.AppendString("\033[90m" + caller.TrimmedPath() + "\033[0m")
	} else {
		enc.AppendString("")
	}
}

func Install() {
	// 控制台编码器配置，仿Spring Boot风格
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:      "T",
		LevelKey:     "L",
		CallerKey:    "C",
		MessageKey:   "M",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder, // 级别彩色
		EncodeTime:   ColorTimeEncoder,                 // 时间绿色
		EncodeCaller: GrayCallerEncoder,                // 文件:行号灰色
	}
	// 字段顺序：时间 级别 文件:行号 : 消息
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 文件编码器配置（无颜色，适合存档）
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	logPath := "logs"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "app.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // 启用压缩
	}

	// 控制台输出（彩色，适合开发）
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	// 文件输出（无色，适合归档）
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), zap.InfoLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(msg string, fields ...any) {
	logger.Sugar().Infof(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Warnf(msg string, fields ...any) {
	logger.Sugar().Warnf(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(msg string, fields ...any) {
	logger.Sugar().Debugf(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorf(msg string, fields ...any) {
	logger.Sugar().Errorf(msg, fields...)
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	return logger
}
