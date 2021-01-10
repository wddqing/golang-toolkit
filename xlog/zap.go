package xlog

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 新建日志, 每30秒更新一次日志等级
func NewLogger(isProduction bool) *zap.Logger {
	logger, _ := newZapLogger(isProduction, os.Stdout)
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return logger
}

func newZapLogger(isProduction bool, output zapcore.WriteSyncer) (*zap.Logger, *zap.AtomicLevel) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
	}

	var encoder zapcore.Encoder
	dyn := zap.NewAtomicLevel()
	if isProduction {
		dyn.SetLevel(zap.InfoLevel)
		encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg) // zapcore.NewJSONEncoder(encCfg)
	} else {
		dyn.SetLevel(zap.DebugLevel)
		encCfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	return zap.New(zapcore.NewCore(encoder, output, dyn), zap.AddCaller()), &dyn
}
