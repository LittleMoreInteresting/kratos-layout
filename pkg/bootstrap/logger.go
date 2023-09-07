package bootstrap

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kratos-layout/pkg/zaplog"
)

func NewLoggerProvider(Id, Name, Version string) log.Logger {
	//l := log.NewStdLogger(os.Stdout)
	l := NewZapLogger(Name)
	return log.With(
		l,
		"service.id", Id,
		"service.name", Name,
		"service.version", Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}

type GormLogger struct {
	l log.Logger
}

func NewGormLogger(Id, Name, Version string) *GormLogger {
	l := NewZapLogger("sql")
	lg := log.With(
		l,
		"service.id", Id,
		"service.name", Name,
		"service.version", Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	return &GormLogger{
		l: lg,
	}
}
func (gl *GormLogger) Printf(s string, v ...interface{}) {
	l := log.NewHelper(log.With(gl.l, "module", "gorm/data/sql"))
	switch {
	case strings.Contains(s, "[info]"):
		l.Infof(s, v...)
	case strings.Contains(s, "[warn]"):
		l.Warnf(s, v...)
	case strings.Contains(s, "[error]"):
		l.Errorf(s, v...)
	default:
		l.Infof(s, v...)
	}
}

// 只有zap 将日志写入文件
func NewZapLogger(prefix string) log.Logger {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	zlogger := zaplog.NewZapLogger(
		encoder,
		prefix,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
	return zlogger
}
