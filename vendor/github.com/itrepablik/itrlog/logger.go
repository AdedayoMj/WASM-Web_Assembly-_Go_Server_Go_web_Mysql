// Package itrlog is the custom logger for Go using Zap and Lumberjack libraries.
package itrlog

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/itrepablik/lumberjack"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
	_logFileName        = "itrlog_"
	_logFolderName      = "logs"
	_logMaxSizeInMB     = 100
	_logMaxAgeInDays    = 0
)

// ITRLogger contains the necessary information to initialize the required parameters for the itrlog package.
type ITRLogger struct {
	MaxSizeInMB, MaxAgeInDays int
	LogFolderName, LogInitial string
	Sugar                     *zap.SugaredLogger
	Base                      *zap.Logger
	mu                        sync.Mutex // ensures atomic writes; protects the following fields
}

// LogTimeFormat formats the event timestamp.
const LogTimeFormat string = "Jan-02-2006 03:04:05 PM"

var s *ITRLogger

// InitLog initialize the zap and lumberjack logger library.
// Example log file: logs/app_name_2020-02-28.log
func InitLog(maxSizeInMB, maxAgeInDays int, logFolderName, logInitial string) *ITRLogger {
	logFile := logFolderName + "/" + logInitial + time.Now().Format("2006-01-02") + ".log"

	// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  logFile,      // log filename
		MaxSize:   maxSizeInMB,  // megabytes
		MaxAge:    maxAgeInDays, // days
		LocalTime: true,         // use the local machine's timezone
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	logger := zap.New(core)
	sugar := logger.Sugar()
	return &ITRLogger{
		MaxSizeInMB:   maxSizeInMB,
		MaxAgeInDays:  maxAgeInDays,
		LogFolderName: logFolderName,
		LogInitial:    logInitial,
		Sugar:         sugar,
		Base:          logger,
	}
}

func init() {
	// Set initial default itrlog required settings.
	s = InitLog(_logMaxSizeInMB, _logMaxAgeInDays, _logFolderName, _logFileName)
}

// SetLogInit sets the custom log requirement to initialize the itr logger.
func SetLogInit(maxSizeInMB, maxAgeInDays int, logFolderName, logFileName string) *ITRLogger {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Set some default minimal requirements
	if maxSizeInMB <= 0 {
		maxSizeInMB = _logMaxSizeInMB
	}
	if maxAgeInDays < 0 {
		maxAgeInDays = _logMaxAgeInDays
	}
	if len(strings.TrimSpace(logFolderName)) == 0 {
		logFolderName = _logFolderName
	}
	if len(strings.TrimSpace(logFileName)) == 0 {
		logFileName = _logFileName
	}

	// Re-configure the itrlog
	s = InitLog(maxSizeInMB, maxAgeInDays, logFolderName, logFileName)
	return s
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	s.log(zap.DebugLevel, "", args, nil)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	s.log(zap.InfoLevel, "", args, nil)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	s.log(zap.WarnLevel, "", args, nil)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	s.log(zap.ErrorLevel, "", args, nil)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	s.log(zap.DPanicLevel, "", args, nil)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	s.log(zap.PanicLevel, "", args, nil)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	s.log(zap.FatalLevel, "", args, nil)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	s.log(zap.DebugLevel, template, args, nil)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	s.log(zap.InfoLevel, template, args, nil)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	s.log(zap.WarnLevel, template, args, nil)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	s.log(zap.ErrorLevel, template, args, nil)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	s.Sugar.DPanicf(template, args)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	s.log(zap.DPanicLevel, template, args, nil)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	s.log(zap.PanicLevel, template, args, nil)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	s.log(zap.DebugLevel, msg, nil, keysAndValues)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	s.log(zap.InfoLevel, msg, nil, keysAndValues)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	s.log(zap.WarnLevel, msg, nil, keysAndValues)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	s.log(zap.ErrorLevel, msg, nil, keysAndValues)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	s.log(zap.DPanicLevel, msg, nil, keysAndValues)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	s.log(zap.PanicLevel, msg, nil, keysAndValues)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	s.log(zap.FatalLevel, msg, nil, keysAndValues)
}

func (s *ITRLogger) log(lvl zapcore.Level, template string, fmtArgs []interface{}, context []interface{}) {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if lvl < zap.DPanicLevel && !s.Base.Core().Enabled(lvl) {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := s.Base.Check(lvl, msg); ce != nil {
		ce.Write(s.sweetenFields(context)...)
	}
}

func (s *ITRLogger) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zap.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			s.Base.DPanic(_oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		s.Base.DPanic(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}