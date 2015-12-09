package gottyclient

import (
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func TestParseURL(t *testing.T) {
	Convey("Testing ParseURL", t, func() {
		Convey("Complete URLs", func() {
			input := "http://test.com:8888/blahblah/blihblih"
			output, err := ParseURL(input)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, input)

			input = "https://test.com:8888/blahblah/blihblih"
			output, err = ParseURL(input)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, input)

			input = "https://test.com:8888"
			output, err = ParseURL(input)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, input)
		})
		Convey("Incomplete URLs", func() {
			input := "test.com:8888/blahblah/blihblih"
			expected := "http://test.com:8888/blahblah/blihblih"
			output, err := ParseURL(input)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, expected)

			input = "test.com"
			expected = "http://test.com"
			output, err = ParseURL(input)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, expected)
		})
	})
}
