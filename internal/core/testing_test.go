package core

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestGoldenIgnoreLines(t *testing.T) {
	original := `
Line1
Line2
Line3
Line4`
	expected := `
Line1
Line4`
	actual := goldenIgnoreLines(original, GoldenIgnoreLine{
		Prefix: "Line2",
		After:  1,
	})
	assert.Equal(t, expected, actual)

	expected2 := `
Line2
Line3`
	actual2 := goldenIgnoreLines(original,
		GoldenIgnoreLine{
			Prefix: "Line1",
			After:  0,
		},
		GoldenIgnoreLine{
			Prefix: "Line4",
			After:  0,
		})
	assert.Equal(t, expected2, actual2)

}
