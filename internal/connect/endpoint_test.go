package connect_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEndpoint struct {
	name    string
	public  bool
	private bool
}

func Test_PreferPublic(t *testing.T) {
	t.Parallel()

	isPublic := func(e fakeEndpoint) bool { return e.public }
	isPrivate := func(e fakeEndpoint) bool { return e.private }

	t.Run("prefers public over private", func(t *testing.T) {
		t.Parallel()
		got, err := connect.PreferPublic([]fakeEndpoint{
			{name: "priv", private: true},
			{name: "pub", public: true},
		}, isPublic, isPrivate)
		require.NoError(t, err)
		assert.Equal(t, "pub", got.name)
	})

	t.Run("falls back to private", func(t *testing.T) {
		t.Parallel()
		got, err := connect.PreferPublic([]fakeEndpoint{
			{name: "priv", private: true},
		}, isPublic, isPrivate)
		require.NoError(t, err)
		assert.Equal(t, "priv", got.name)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		_, err := connect.PreferPublic(nil, isPublic, isPrivate)
		require.EqualError(t, err, "no endpoint found")
	})

	t.Run("neither public nor private", func(t *testing.T) {
		t.Parallel()
		_, err := connect.PreferPublic([]fakeEndpoint{{name: "other"}}, isPublic, isPrivate)
		require.EqualError(t, err, "no endpoint found")
	})
}
