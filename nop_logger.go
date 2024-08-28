package apilog

import "time"

// NewNop returns a no-op Logger. Do nothing and never writes out any logs.
func NewNop() Logger {
	return &nopLogger{}
}

type nopLogger struct{}

func (n nopLogger) Init(_ time.Duration)            {}
func (n nopLogger) Flush(_ time.Duration)           {}
func (n nopLogger) With(_ ...Log) Logger            { return n }
func (n nopLogger) Group(_ string, _ ...Log) Logger { return n }
func (n nopLogger) Dbg(_ string, _ ...Log)          {}
func (n nopLogger) Inf(_ string, _ ...Log)          {}
func (n nopLogger) Wrn(_ string, _ ...Log)          {}
func (n nopLogger) Err(_ string, _ ...Log)          {}
