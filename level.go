package apilog

import "strings"

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel most verbose logs, and are usually disabled in production.
	DebugLevel Level = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

// ParseLevel parses a level based on the lower-case representation of the log
// level.
func ParseLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning", "warn":
		return WarnLevel
	case "error", "err":
		return ErrorLevel
	}
	return -1
}
