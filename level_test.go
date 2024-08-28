package apilog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		name   string
		sample string
		expect Level
	}{
		{
			name:   "Debug",
			sample: "debug",
			expect: DebugLevel,
		},
		{
			name:   "Info",
			sample: "info",
			expect: InfoLevel,
		},
		{
			name:   "Warn",
			sample: "warn",
			expect: WarnLevel,
		},
		{
			name:   "Warning",
			sample: "warning",
			expect: WarnLevel,
		},
		{
			name:   "Err",
			sample: "err",
			expect: ErrorLevel,
		},
		{
			name:   "Error",
			sample: "error",
			expect: ErrorLevel,
		},
		{
			name:   "Unrecognized",
			sample: "hello",
			expect: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, ParseLevel(tc.sample))
		})
	}
}
