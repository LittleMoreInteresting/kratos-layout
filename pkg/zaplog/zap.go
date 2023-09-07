package zaplog

import (
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

// ZapLogger is a logger impl.
type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, prefix string, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	rotateLogger, err := rotatelogs.New("./log/"+prefix+"_%Y%m%d.log", rotatelogs.WithMaxAge(time.Hour*24*30))
	if err != nil {
		log.Fatalf("logger.Setup err: %v", err)
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),    //输出到控制台上
			zapcore.AddSync(rotateLogger), //输出到文本中
		), level)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	}
	return nil
}
