package main

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RunMode int

const (
	Development RunMode = 0
	Production  RunMode = 1
)

var elogger *zap.SugaredLogger
var inited bool = false

func GetLogger() (*zap.SugaredLogger, error) {
	if inited {
		return elogger, nil
	} else {
		return nil, errors.New("not initialized")
	}
}

// mode: 0:development, 1:production
func InitLog(mode RunMode, file_path string, size int, backups int, age int, to_console bool) {
	if inited {
		return
	}

	hook := lumberjack.Logger{
		Filename:   file_path,
		MaxSize:    size,    // MB
		MaxBackups: backups, // backup files
		MaxAge:     age,     // days
		Compress:   true,
	}

	encoder_config := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "v",
		NameKey:        "logger",
		CallerKey:      "l",
		MessageKey:     "m",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"),
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// encoder_config := zap.NewProductionEncoderConfig()
	// encoder_config.EncodeTime = zapcore.ISO8601TimeEncoder

	atomicLevel := zap.NewAtomicLevel()
	if mode == Development {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	var wsync zapcore.WriteSyncer
	if to_console {
		wsync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		wsync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder_config),
		// zapcore.NewConsoleEncoder(encoder_config),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		wsync,
		atomicLevel,
	)

	// caller := zap.AddCaller()
	// development := zap.Development()
	// filed := zap.Fields(zap.String("serviceName1", "serviceName2"))
	// logger := zap.New(core, caller, development, filed)
	// elogger = zap.New(core, caller, development).Sugar()

	if mode == Development {
		elogger = zap.New(core, zap.AddCaller(), zap.Development()).Sugar()
	} else {
		elogger = zap.New(core, zap.AddCaller()).Sugar()
	}

	inited = true
}
