package core_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/core"
)

func TestOneOf(t *testing.T) {
	a := &core.ArgSpec{
		Name:       "Argument A",
		OneOfGroup: "ab group",
	}
	b := &core.ArgSpec{
		Name:       "Argument B",
		OneOfGroup: "ab group",
	}
	c := &core.ArgSpec{
		Name:       "Argument C",
		OneOfGroup: "",
	}
	d := &core.ArgSpec{
		Name:       "Argument D",
		OneOfGroup: "",
	}
	e := &core.ArgSpec{
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
