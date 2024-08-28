package apilog

import (
	"io"
	"time"
)

// Writer unified log writer that's responsible where the log from Logger
// should be written to.
type Writer interface {
	// Writer return where and how the implementer should write the logs.
	Writer() io.Writer
	// Output define the Output type.
	Output() Output
	// Level define the logs level.
	Level() Level
	// Wait in case the implementer need some delay or preparation before can write any logs.
	Wait(dur time.Duration)
	// Flush any necessary clean up task that will be run by Producer at the last order.
	Flush(dur time.Duration)
}

// Output define currently supported target output for logging.
type Output int8

const (
	CONSOLE  Output = iota // CONSOLE target log output to console/terminal
	NEWRELIC               // NEWRELIC target log output directly to new relic via newrelic client sdk
	FILE                   // FILE target log output to local file
)
