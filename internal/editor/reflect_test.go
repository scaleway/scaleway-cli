package editor

import (
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_valueMapper(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 string
		Arg2 int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Equal(t, src.Arg1, dest.Arg1)
	assert.Equal(t, src.Arg2, dest.Arg2)
}

func Test_valueMapperInvalidType(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 string
		Arg2 bool
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Equal(t, src.Arg1, dest.Arg1)
	assert.Zero(t, dest.Arg2)
}

func Test_valueMapperDifferentFields(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg3 string
		Arg4 bool
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.Zero(t, dest.Arg3)
	assert.Zero(t, dest.Arg4)
}

func Test_valueMapperPointers(t *testing.T) {
	src := struct {
		Arg1 string
		Arg2 int
	}{"1", 1}
	dest := struct {
		Arg1 *string
		Arg2 *int
	}{}

	valueMapper(reflect.ValueOf(&dest), reflect.ValueOf(&src))
	assert.NotNil(t, dest.Arg1)
	assert.EqualValues(t, src.Arg1, *dest.Arg1)
	assert.NotNil(t, dest.Arg2)
	assert.EqualValues(t, src.Arg2, *dest.Arg2)
}
