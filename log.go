package apilog

// Log object that holds data for each field inserted to each log message. How
// Logger implementer is treating this object should read the field typ and
// follow the guideline from Type and each of the supported types.
type Log struct {
	typ Type
	key string
	str string
	num int
	flt float64
	b   bool
	any any
	err error
}

// Type indicates how the Logger implementer should treat each
// Log.
type Type uint8

const (
	// StringType use field str string of Log as the value.
	StringType Type = iota
	// NumType use field num int of Log as the value.
	NumType
	// FloatType use field flt float64 of Log as the value.
	FloatType
	// BoolType use field b bool of Log as the value.
	BoolType
	// AnyType use field any interface of Log as the value.
	AnyType
	// ErrorType use field err from error interface of Log as the value.
	ErrorType
)

// String constructs a Log with the given key and value. This set the type to
// StringType.
func String(k, v string) Log {
	return Log{typ: StringType, key: k, str: v}
}

// Num constructs a Log with the given key and value. This set the type to
// NumType.
func Num(k string, num int) Log {
	return Log{typ: NumType, key: k, num: num}
}

// Float constructs a Log with the given key and value. This set the type to
// FloatType.
func Float(k string, f float64) Log {
	return Log{typ: FloatType, key: k, flt: f}
}

// Bool constructs a Log with the given key and value. This set the type to
// BoolType.
func Bool(k string, b bool) Log {
	return Log{typ: BoolType, key: k, b: b}
}

// Any constructs a Log with the given key and value. This set the type to
// AnyType.
func Any(k string, any any) Log {
	return Log{typ: AnyType, key: k, any: any}
}

// Error constructs a Log with the given err value and 'error' as the key. This
// set the type to ErrorType.
func Error(err error) Log {
	return Log{typ: ErrorType, key: "error", err: err}
}
