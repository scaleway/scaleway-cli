package core

import (
	"net/http"
	"time"
)

const defaultRetryInterval = 1 * time.Second

type retryableHTTPTransport struct {
	transport http.RoundTripper
}

func (r *retryableHTTPTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	res, err := r.transport.RoundTrip(request)
	if err == nil && res.StatusCode == http.StatusTooManyRequests {
		time.Sleep(defaultRetryInterval)
		return r.RoundTrip(request)
	}
	return res, err
}
