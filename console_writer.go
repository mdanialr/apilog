package apilog

import (
	"io"
	"os"
	"time"
)

// NewConsoleWriter return Writer implementer that write logs to os.Stdout and
// set given lvl as the log Level.
func NewConsoleWriter(lvl Level) Writer {
	return &consoleOutput{lvl: lvl}
}

type consoleOutput struct {
	lvl Level
}

func (c *consoleOutput) Writer() io.Writer     { return os.Stdout }
func (c *consoleOutput) Output() Output        { return CONSOLE }
func (c *consoleOutput) Level() Level          { return c.lvl }
func (c *consoleOutput) Wait(_ time.Duration)  {}
func (c *consoleOutput) Flush(_ time.Duration) {}
