package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.SugaredLogger
	once   sync.Once
)

// Init initializes the global logger. In debug mode it uses a human-readable
// console encoder on stderr at Debug level. In production mode it uses a JSON
// encoder on stderr at Info level.
func Init(debug bool) {
	once.Do(func() {
		global = newLogger(debug)
	})
}

// L returns the global sugared logger. Safe to call before Init (returns a
// no-op logger).
func L() *zap.SugaredLogger {
	if global == nil {
		return zap.NewNop().Sugar()
	}
	return global
}

// Sync flushes any buffered log entries.
func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func newLogger(debug bool) *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zapcore.Level
	var encoder zapcore.Encoder

	if debug {
		level = zapcore.DebugLevel
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		level = zapcore.InfoLevel
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		level,
	)

	opts := []zap.Option{zap.AddCaller()}
	if debug {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, opts...).Sugar()
}
