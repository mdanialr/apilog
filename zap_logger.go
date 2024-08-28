package apilog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger return Logger implementer that use zap as the backend.
func NewZapLogger(wr ...Writer) Logger {
	// set to singleton
	singletonLogger = &zapLogger{wr: wr}
	return singletonLogger
}

type zapLogger struct {
	log *zap.Logger
	wr  []Writer
}

func (z *zapLogger) clone() *zapLogger {
	c := *z
	return &c
}

func (z *zapLogger) Init(dur time.Duration) {
	var cores []zapcore.Core
	// setup common zap json encoder
	jsonEnc := zap.NewProductionEncoderConfig()
	jsonEnc.EncodeTime = zapcore.RFC3339TimeEncoder
	jsonEnc.EncodeLevel = zapcore.CapitalLevelEncoder
	jsonEnc.TimeKey = "time"

	for _, w := range z.wr {
		switch w.Output() {
		case CONSOLE:
			encCnf := zap.NewDevelopmentConfig().EncoderConfig
			encCnf.EncodeLevel = zapcore.CapitalColorLevelEncoder
			enc := zapcore.NewConsoleEncoder(encCnf)
			core := zapcore.NewCore(enc, zapcore.AddSync(w.Writer()), toZapLevel(w.Level()))
			cores = append(cores, core)

		case FILE, NEWRELIC:
			enc := zapcore.NewJSONEncoder(jsonEnc)
			core := zapcore.NewCore(enc, zapcore.AddSync(w.Writer()), toZapLevel(w.Level()))
			cores = append(cores, core)
		}
		w.Wait(dur)
	}
	z.log = zap.New(zapcore.NewTee(cores...))
}

func (z *zapLogger) Flush(dur time.Duration) {
	for _, w := range z.wr {
		w.Flush(dur)
	}
	_ = z.log.Sync()
}

func (z *zapLogger) With(pr ...Log) Logger {
	if len(pr) == 0 {
		return z
	}
	mutex.Lock()
	defer mutex.Unlock()

	// clone it, so on every With method call does not affect the parent logger
	clone := z.clone()
	clone.log = clone.log.With(toZapFields(pr)...)

	// then reassign to singleton
	singletonLogger = clone

	return clone
}

func (z *zapLogger) Group(key string, pr ...Log) Logger {
	if len(pr) == 0 || key == "" {
		return z
	}
	mutex.Lock()
	defer mutex.Unlock()

	// clone it, so on every Group method call does not affect the parent logger
	clone := z.clone()
	// for zap, the trick is to use Any instead of Namespace, because as the docs said
	//  it will eat any subsequent context data and treat it as their fields
	//   https://pkg.go.dev/go.uber.org/zap#Namespace ('... All subsequent fields will be added to the new namespace.')
	//
	// Example: 'ctx' is the subsequent context data
	//  zap.Namespace() => {"repo_layer":{"id":"1234567890","ctx":{"agent":"007"}}}
	//  zap.Any() => {"repo_layer":{"id":"1234567890"},"ctx":{"agent":"007"}}
	//
	//  notice that 'ctx' is embedded as 'repo_layer' field when using Namespace, but
	//   it's properly wrapped as intended when using Any
	clone.log = clone.log.With(zap.Any(key, toZapFields(pr)))

	// then reassign to singleton
	singletonLogger = clone

	return clone
}

func (z *zapLogger) Dbg(msg string, pr ...Log) {
	if len(pr) > 0 {
		z.log.Debug(msg, toZapFields(pr)...)
		return
	}
	z.log.Debug(msg)
}

func (z *zapLogger) Inf(msg string, pr ...Log) {
	if len(pr) > 0 {
		z.log.Info(msg, toZapFields(pr)...)
		return
	}
	z.log.Info(msg)
}

func (z *zapLogger) Wrn(msg string, pr ...Log) {
	if len(pr) > 0 {
		z.log.Warn(msg, toZapFields(pr)...)
		return
	}
	z.log.Warn(msg)
}

func (z *zapLogger) Err(msg string, pr ...Log) {
	if len(pr) > 0 {
		z.log.Error(msg, toZapFields(pr)...)
		return
	}
	z.log.Error(msg)
}

// toZapLevel transform log Level to zap level.
func toZapLevel(lvl Level) zapcore.Level {
	switch lvl {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	}
	return zapcore.InvalidLevel
}

// toZapFields transform Log to zap field.
func toZapFields(pr []Log) []zapcore.Field {
	var fields []zapcore.Field
	for _, p := range pr {
		switch p.typ {
		case StringType:
			fields = append(fields, zap.String(p.key, p.str))
		case NumType:
			fields = append(fields, zap.Int(p.key, p.num))
		case FloatType:
			fields = append(fields, zap.Float64(p.key, p.flt))
		case BoolType:
			fields = append(fields, zap.Bool(p.key, p.b))
		case AnyType:
			fields = append(fields, zap.Any(p.key, p.any))
		case ErrorType:
			fields = append(fields, zap.NamedError(p.key, p.err))
		}
	}
	return fields
}
