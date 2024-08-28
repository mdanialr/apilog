package apilog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsoleWriter(t *testing.T) {
	cns := NewConsoleWriter(DebugLevel)

	assert.Equal(t, DebugLevel, cns.Level())
	assert.Equal(t, os.Stdout, cns.Writer())
	assert.Equal(t, CONSOLE, cns.Output())

	cns.Wait(1)  // do nothing
	cns.Flush(1) // do nothing
}
