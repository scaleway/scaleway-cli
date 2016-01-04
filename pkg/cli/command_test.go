package cli

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand_Name(t *testing.T) {
	Convey("Testing Command.Name()", t, func() {
		command := Command{
			UsageLine: "top [OPTIONS] SERVER",
		}
		So(command.Name(), ShouldEqual, "top")

		command = Command{
			UsageLine: "top",
		}
		So(command.Name(), ShouldEqual, "top")
	})
}
