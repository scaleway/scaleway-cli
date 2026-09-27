package core

import (
	"net/http"
	"testing"

	"github.com/scaleway/scaleway-sdk-go/vcr"
	"github.com/stretchr/testify/assert"
)


func newHTTPRecorder(t *testing.T, folder string, update bool) (*http.Client, func(), error) {
	t.Helper()

	r, err := vcr.NewHTTPRecorder(t, folder, update, &SocketPassthroughTransport{})
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{Transport: &retryableHTTPTransport{transport: r}}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
