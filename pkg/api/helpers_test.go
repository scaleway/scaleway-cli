package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type VolumesFromSizeCase struct {
	name  string
	input struct {
		rootVolumeSize, targeSize, perVolumeMaxSize uint64
	}
	output string
}

func TestVolumesFromSize(t *testing.T) {
	tests := []VolumesFromSizeCase{
		{
			name: "200G 200G 200G",
			input: struct{ rootVolumeSize, targeSize, perVolumeMaxSize uint64 }{
				200 * Giga, 600 * Giga, 200 * Giga,
			},
			output: "200G 200G",
		},
		{
			name: "200G 200G 100G",
			input: struct{ rootVolumeSize, targeSize, perVolumeMaxSize uint64 }{
				200 * Giga, 500 * Giga, 200 * Giga,
			},
			output: "200G 100G",
		},
		{
			name: "25G",
			input: struct{ rootVolumeSize, targeSize, perVolumeMaxSize uint64 }{
				25 * Giga, 25 * Giga, 200 * Giga,
			},
			output: "",
		},
		{
			name: "100G 150G",
			input: struct{ rootVolumeSize, targeSize, perVolumeMaxSize uint64 }{
				100 * Giga, 250 * Giga, 200 * Giga,
			},
			output: "150G",
		},
		{
			name: "200G 50G 50G 50G 50G",
			input: struct{ rootVolumeSize, targeSize, perVolumeMaxSize uint64 }{
				200 * Giga, 400 * Giga, 50 * Giga,
			},
			output: "50G 50G 50G 50G",
		},
	}
	for _, test := range tests {
		Convey("Testing VolumesFromSize with expected "+test.name, t, func(c C) {
			output := VolumesFromSize(test.input.rootVolumeSize, test.input.targeSize, test.input.perVolumeMaxSize)
			c.So(output, ShouldEqual, test.output)
		})
	}

}
