package apilog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewNewrelicWriter return Writer implementer that ingest logs directly to
// newrelic server by given Config and set given Level as the log level.
func NewNewrelicWriter(lvl Level, cnf *Config) Writer {
	if cnf == nil {
		cnf = &Config{}
	}

	nr, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cnf.nr.name),
		newrelic.ConfigLicense(cnf.nr.license),
		newrelic.ConfigInfoLogger(os.Stdout),
	)
	if err != nil {
		panic(errors.New("failed to init newrelic writer: " + err.Error()))
	}
	return &newrelicOutput{lvl: lvl, nr: nr}
}

type newrelicOutput struct {
	nr  *newrelic.Application
	lvl Level
}

// Write implement io.Writer by passing the data to newrelic app.
func (n *newrelicOutput) Write(p []byte) (int, error) {
	msg := bytes.TrimSpace(p)
	n.nr.RecordLog(newrelic.LogData{Message: string(msg)})
	return len(p), nil
}
func (n *newrelicOutput) Writer() io.Writer       { return n }
func (n *newrelicOutput) Output() Output          { return NEWRELIC }
func (n *newrelicOutput) Level() Level            { return n.lvl }
func (n *newrelicOutput) Wait(dur time.Duration)  { _ = n.nr.WaitForConnection(dur) }
func (n *newrelicOutput) Flush(dur time.Duration) { n.nr.Shutdown(dur) }
