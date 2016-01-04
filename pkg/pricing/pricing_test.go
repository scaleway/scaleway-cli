package pricing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetByPath(t *testing.T) {
	Convey("Testing GetByPath", t, func() {
		object := CurrentPricing.GetByPath("/compute/c1/run")
		So(object, ShouldNotBeNil)
		So(object.Path, ShouldEqual, "/compute/c1/run")

		object = CurrentPricing.GetByPath("/ip/dynamic")
		So(object, ShouldNotBeNil)
		So(object.Path, ShouldEqual, "/ip/dynamic")

		object = CurrentPricing.GetByPath("/dontexists")
		So(object, ShouldBeNil)
	})
}

func TestGetByIdentifier(t *testing.T) {
	Convey("Testing GetByIdentifier", t, func() {
		object := CurrentPricing.GetByIdentifier("aaaaaaaa-aaaa-4aaa-8aaa-111111111112")
		So(object, ShouldNotBeNil)
		So(object.Path, ShouldEqual, "/compute/c1/run")

		object = CurrentPricing.GetByIdentifier("dontexists")
		So(object, ShouldBeNil)
	})
}
