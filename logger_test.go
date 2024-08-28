package apilog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNop(t *testing.T) {
	t.Run("Should have nop Log type", func(t *testing.T) {
		nl := NewNop()
		assert.NotNil(t, nl)
		assert.Equal(t, &nopLogger{}, nl)
		assert.IsType(t, &nopLogger{}, nl)
	})

	t.Run("Should do nothing", func(t *testing.T) {
		nl := NewNop()
		assert.NotNil(t, nl)
		assert.NotNil(t, nl.With())
		assert.NotNil(t, nl.Group("hello"))

		// just run it, since its just do literary nothing
		nl.Init(time.Microsecond)
		nl.Flush(time.Microsecond)
		nl.Dbg("")
		nl.Inf("")
		nl.Wrn("")
		nl.Err("")
	})
}
