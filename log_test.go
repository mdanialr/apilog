package apilog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	str := String("key", "val")
	assert.Equal(t, StringType, str.typ)
	assert.Equal(t, "key", str.key)
	assert.Equal(t, "val", str.str)

	num := Num("num", 11)
	assert.Equal(t, NumType, num.typ)
	assert.Equal(t, "num", num.key)
	assert.Equal(t, 11, num.num)

	fl := Float("float", 1.1)
	assert.Equal(t, FloatType, fl.typ)
	assert.Equal(t, "float", fl.key)
	assert.Equal(t, 1.1, fl.flt)

	b := Bool("boolean", true)
	assert.Equal(t, BoolType, b.typ)
	assert.Equal(t, "boolean", b.key)
	assert.Equal(t, true, b.b)

	m := make(map[string]any)
	m["object"] = "value"
	an := Any("anything", m)
	assert.Equal(t, AnyType, an.typ)
	assert.Equal(t, "anything", an.key)
	assert.Equal(t, m, an.any)

	er := errors.New("oops")
	err := Error(er)
	assert.Equal(t, ErrorType, err.typ)
	assert.Equal(t, "error", err.key)
	assert.Equal(t, er, err.err)
	assert.Equal(t, "oops", err.err.Error())
}
