package core

import (
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
