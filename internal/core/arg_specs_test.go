package core

import (
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
)

func TestOneOf(t *testing.T) {
	a := &ArgSpec{
		Name:       "Argument A",
		OneOfGroup: "ab group",
	}
	b := &ArgSpec{
		Name:       "Argument B",
		OneOfGroup: "ab group",
	}
	c := &ArgSpec{
		Name:       "Argument C",
		OneOfGroup: "",
	}
	d := &ArgSpec{
		Name:       "Argument D",
		OneOfGroup: "",
	}
	e := &ArgSpec{
		Name: "Argument E",
	}
	assert.True(t, a.ConflictWith(b))
	assert.True(t, b.ConflictWith(a))
	assert.False(t, c.ConflictWith(d))
	assert.False(t, d.ConflictWith(c))
	assert.False(t, a.ConflictWith(c))
	assert.False(t, d.ConflictWith(b))
	assert.False(t, a.ConflictWith(e))
	assert.False(t, a.ConflictWith(c))
	assert.False(t, e.ConflictWith(e))
}

func TestArgSpecGetArgsTypeField(t *testing.T) {
	data := struct {
		Field       string
		FieldStruct struct {
			NestedField int
		}
		FieldSlice []float32
		FieldMap   map[string]bool
	}{}
	dataType := reflect.TypeOf(data)

	fieldSpec := ArgSpec{Name: "field"}
	typ, err := fieldSpec.GetArgsTypeField(dataType)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf("string"), typ, "%s is not string", typ.Name())

	fieldSpec = ArgSpec{Name: "field-struct.nested-field"}
	typ, err = fieldSpec.GetArgsTypeField(dataType)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf(int(1)), typ, "%s is not int", typ.Name())

	fieldSpec = ArgSpec{Name: "field-slice.{index}"}
	typ, err = fieldSpec.GetArgsTypeField(dataType)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf(float32(1)), typ, "%s is not float32", typ.Name())

	fieldSpec = ArgSpec{Name: "field-map.{key}"}
	typ, err = fieldSpec.GetArgsTypeField(dataType)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf(true), typ, "%s is not bool", typ.Name())
}
