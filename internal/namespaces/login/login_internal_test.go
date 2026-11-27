package login

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLoginURL(t *testing.T) {
	callbackURL := "http://localhost:1234/callback"

	t.Run("without organization", func(t *testing.T) {
		assert.Equal(
			t,
			"https://account.scaleway.com/authenticate?redirectToUrl=http%3A%2F%2Flocalhost%3A1234%2Fcallback",
			buildLoginURL(callbackURL, ""),
		)
	})

	t.Run("with organization", func(t *testing.T) {
		loginURL, err := url.Parse(
			buildLoginURL(callbackURL, "11111111-1111-1111-1111-111111111111"),
		)
		assert.NoError(t, err)
		assert.Equal(t, "account.scaleway.com", loginURL.Host)
		assert.Equal(t, "/member-login", loginURL.Path)
		assert.Equal(
			t,
			"11111111-1111-1111-1111-111111111111",
			loginURL.Query().Get("organizationIdOrAlias"),
		)
		assert.Equal(
			t,
			buildLoginURL(callbackURL, ""),
			loginURL.Query().Get("redirectToUrl"),
		)
	})
}
