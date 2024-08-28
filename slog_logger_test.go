package apilog

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToSlogLevel(t *testing.T) {
	testCases := []struct {
		name   string
		sample Level
		expect slog.Level
	}{
		{
			name:   "Debug level",
			sample: DebugLevel,
			expect: slog.LevelDebug,
		},
		{
			name:   "Info level",
			sample: InfoLevel,
			expect: slog.LevelInfo,
		},
		{
			name:   "Warn level",
			sample: WarnLevel,
			expect: slog.LevelWarn,
		},
		{
			name:   "Error level",
			sample: ErrorLevel,
			expect: slog.LevelError,
		},
		{
			name:   "Unrecognized level",
			sample: -1,
			expect: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, toSlogLevel(tc.sample))
		})
	}
}

func TestToSlogAttr(t *testing.T) {
	testCases := []struct {
		name   string
		sample Log
		expect []any
	}{
		{
			name:   "String attribute",
			sample: String("hello", "world"),
			expect: []any{slog.String("hello", "world")},
		},
		{
			name:   "Numeric (int) attribute",
			sample: Num("number", 9),
			expect: []any{slog.Int("number", 9)},
		},
		{
			name:   "Float (decimal number) attribute",
			sample: Float("scale", 0.2),
			expect: []any{slog.Float64("scale", 0.2)},
		},
		{
			name:   "Boolean attribute",
			sample: Bool("is_test", true),
			expect: []any{slog.Bool("is_test", true)},
		},
		{
			name:   "Any (interface) attribute",
			sample: Any("anything", map[string]string{"hi": "hi"}),
			expect: []any{slog.Any("anything", map[string]string{"hi": "hi"})},
		},
		{
			name:   "Error (any) attribute",
			sample: Error(errors.New("oops")),
			expect: []any{slog.Any("error", errors.New("oops"))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, toSlogAttr([]Log{tc.sample}))
		})
	}
}

func TestMultiSlog(t *testing.T) {
	var ms = new(multiSlog)
	var buf = new(bytes.Buffer)
	sl := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	ms.loggers = append(ms.loggers, sl)
	ms = ms.With(slog.String("hello", "world"))
	ms.Debug("debug")
	ms.Info("info")
	ms.Warn("warning")
	ms.Error("error", slog.String("key", "value"))
	msg := strings.Split(strings.TrimSpace(buf.String()), "\n")

	require.Len(t, msg, 4)
	// assert the log message
	dbg := `"level":"DEBUG","msg":"debug","hello":"world"`
	assert.Contains(t, msg[0], dbg)
	inf := `"level":"INFO","msg":"info","hello":"world"`
	assert.Contains(t, msg[1], inf)
	wrn := `"level":"WARN","msg":"warning","hello":"world"`
	assert.Contains(t, msg[2], wrn)
	err := `"level":"ERROR","msg":"error","hello":"world","key":"value"`
	assert.Contains(t, msg[3], err)
}

func TestNewSlogLogger(t *testing.T) {
	t.Run("Console Writer type", func(t *testing.T) {
		// setup
		writer, obs := NewObserverWriter(DebugLevel, CONSOLE)
		wr := NewSlogLogger(writer)
		wr.Init(time.Microsecond)

		// start logging message
		wr = wr.With(String("hello", "world"))
		wr.Dbg("debug log")
		wr.Inf("info log")
		wr.Wrn("warning log")
		wr.Err("error log")

		// just assert the logs len, since the slog console handler is not json string
		require.Equal(t, 4, obs.Len())
		wr.Flush(time.Microsecond)
	})

	t.Run("File or Newrelic Writer type", func(t *testing.T) {
		// setup
		writer, obs := NewObserverWriter(DebugLevel, FILE)
		wr := NewSlogLogger(writer)
		wr.Init(time.Microsecond)

		// start logging message
		wr = wr.With(String("hello", "world"))
		wr = wr.With() // just to increase code coverage
		wr.Dbg("debug log", Num("number", 11))
		wr.Inf("info log", Bool("ok", true))
		wr.Wrn("warning log", Float("scale", 1.2))
		wr.Err("error log", Error(errors.New("oops")))

		require.Equal(t, 4, obs.Len())
		// assert debug log
		dbg := obs.All()[0]
		assert.True(t, dbg.EqualLevel(DebugLevel))
		assert.True(t, dbg.EqualMsg("debug log"))
		assert.Equal(t, "world", dbg.Get("hello"))
		assert.Equal(t, 11, int(dbg.Get("number").(float64)))

		logs := obs.TakeAll()
		// now All() should cause panic
		require.Panics(t, func() {
			_ = obs.All()[0]
		})

		// assert info log
		inf := logs[1]
		assert.True(t, inf.EqualLevel(InfoLevel))
		assert.True(t, inf.ContainMsg("info log"))
		assert.Equal(t, "world", inf.Get("hello"))
		assert.Equal(t, true, inf.Get("ok"))
		// assert warning log
		wrn := logs[2]
		assert.True(t, wrn.EqualLevel(WarnLevel))
		assert.True(t, wrn.EqualMsg("warning log"))
		assert.Equal(t, 1.2, wrn.Get("scale"))
		// assert error log
		err := logs[3]
		assert.True(t, err.EqualLevel(ErrorLevel))
		assert.True(t, err.ContainMsg("error log"))
		assert.Equal(t, "world", err.Get("hello"))
		assert.Equal(t, "oops", err.Get("error"))

		wr.Flush(time.Microsecond)
	})
}
