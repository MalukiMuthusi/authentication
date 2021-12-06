package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zlog *zap.Logger
var enableDebugMessage bool

var zLock = sync.RWMutex{}

// TODO: Extract to a share library that is a separate git repository
func Setup(dbg bool) (err error) {
	zLock.Lock()
	defer zLock.Unlock()
	enableDebugMessage = dbg
	if enableDebugMessage {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableStacktrace = true
		cfg.DisableCaller = true
		zlog, err = cfg.Build(zap.AddStacktrace(zap.ErrorLevel))
	} else {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableCaller = true
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		zlog, err = cfg.Build()
	}
	if err == nil {
		zlog.Sync()
	}
	return
}

//Info logs with info level
func Info(m string, fields ...zapcore.Field) {
	zlog.Info(m, fields...)
}

//Warn logs with warn level
func Warn(m string, fields ...zapcore.Field) {
	zlog.Warn(m, fields...)
}

// Debug logs with debug level
func Debug(m string, fields ...zapcore.Field) {
	zlog.Debug(m, fields...)
}

// Error logs with error level
func Error(m string, fields ...zapcore.Field) {
	zlog.Error(m, fields...)
}

// Fatalf logs with fatal level and calls os.Exit(1)
func Fatalf(m string, fields ...zapcore.Field) {
	zlog.Fatal(m, fields...)
}
