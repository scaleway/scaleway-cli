package core

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestShell_OptionToArgSpecName(t *testing.T) {
	tt := []struct {
		Option      string
		ArgSpecName string
	}{
		{
			"additional-volumes.0=hello",
			"additional-volumes.{index}",
		},
		{
			"pools.0.kubelet-args.",
			"pools.{index}.kubelet-args.{key}",
		},
	}
	for _, test := range tt {
		assert.Equal(t, optionToArgSpecName(test.Option), test.ArgSpecName)
	}
}

func TestShell_isOption(t *testing.T) {
	tt := []struct {
		Arg      string
		IsOption bool
	}{
		{
			"image=",
			true,
		},
		{
			"pools.0.autoscaling=",
			true,
		},
		{
			"pools.0.kubelet-args.",
			true,
		},
	}
	for _, test := range tt {
		assert.Equal(t, argIsOption(test.Arg), test.IsOption, "%s option value is wrong", test.Arg)
	}
}
