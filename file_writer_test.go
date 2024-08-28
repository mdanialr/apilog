package apilog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestSetupLumberjack(t *testing.T) {
	t.Run("Given empty value config should just use the default value instead", func(t *testing.T) {
		ll := setupLumberjack(&FileConfig{})
		assert.Equal(t, "./logs/app.log", ll.Filename)
		assert.Equal(t, 150, ll.MaxSize)
		assert.Equal(t, 28, ll.MaxAge)
		assert.Equal(t, 7, ll.MaxBackups)
	})

	t.Run("Given config with certain value should use that value instead", func(t *testing.T) {
		fc := &FileConfig{
			path: "./app-logs/app.log",
			size: 200,
			age:  25,
			num:  10,
		}
		ll := setupLumberjack(fc)
		assert.Equal(t, "./app-logs/app.log", ll.Filename)
		assert.Equal(t, 200, ll.MaxSize)
		assert.Equal(t, 25, ll.MaxAge)
		assert.Equal(t, 10, ll.MaxBackups)
	})
}

func TestNewFileWriter(t *testing.T) {
	t.Run("Should return the expected value in each Writer implementation", func(t *testing.T) {
		wr := NewFileWriter(ErrorLevel, nil)
		assert.IsType(t, &lumberjack.Logger{}, wr.Writer())
		assert.Equal(t, FILE, wr.Output())
		assert.Equal(t, ErrorLevel, wr.Level())

		// just run
		wr.Wait(-1)
		wr.Flush(-1)
	})
}
