package api

import (
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func TestNewScalewayAPI(t *testing.T) {
	Convey("Testing NewScalewayAPI()", t, func() {
		api, err := NewScalewayAPI("http://api-endpoint.com", "http://account-endpoint.com", "my-organization", "my-token")
		So(err, ShouldBeNil)
		So(api, ShouldNotBeNil)
		So(api.ComputeAPI, ShouldEqual, "http://api-endpoint.com")
		So(api.AccountAPI, ShouldEqual, "http://account-endpoint.com")
		So(api.Token, ShouldEqual, "my-token")
		So(api.Organization, ShouldEqual, "my-organization")
		So(api.Cache, ShouldNotBeNil)
		So(api.client, ShouldNotBeNil)
		So(api.anonuuid, ShouldNotBeNil)
	})
}
