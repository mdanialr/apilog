package apilog

import "time"

// Logger unified front-end to log.
type Logger interface {
	// Init do initialization and should be called at very first before
	// any other functions.
	Init(dur time.Duration)
	// Flush do clean up task that make sure any leftover logs written properly
	// to the destination by the Log therefor this should be called at very
	// last after other functions e.g. when gracefully shutting down server.
	Flush(dur time.Duration)
	// With add given Log(s) as structured context.
	With(pr ...Log) Logger
	// Group create new group or namespace with given key and Log(s) as the
	// content. In JSON, this is like creating new object using given key and
	// Log(s) as the fields.
	//
	// Ref(s):
	//  - https://pkg.go.dev/go.uber.org/zap#Namespace
	//  - https://pkg.go.dev/golang.org/x/exp/slog#Group
	Group(key string, pr ...Log) Logger
	// Dbg logs a message at DebugLevel.
	Dbg(msg string, pr ...Log)
	// Inf logs a message at InfoLevel.
	Inf(msg string, pr ...Log)
	// Wrn logs a message at WarnLevel.
	Wrn(msg string, pr ...Log)
	// Err logs a message at ErrorLevel.
	Err(msg string, pr ...Log)
}
