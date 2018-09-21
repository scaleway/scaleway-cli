package cli

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecurity_groups(t *testing.T) {
	Convey("Testing Security groups valid policy", t, func() {
		So(isValidPolicy("accept"), ShouldEqual, true)
		So(isValidPolicy("notvalid"), ShouldEqual, false)
	})
}
