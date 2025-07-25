package core

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/scaleway/scaleway-sdk-go/logger"
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

func (r *retryableHTTPTransport) SetInsecureTransport() {
	transportClient, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		logger.Warningf("cli: cannot use insecure mode with DefaultTransport of type %T", http.DefaultTransport)
		return
	}
	if transportClient.TLSClientConfig == nil {
		transportClient.TLSClientConfig = &tls.Config{}
	}
	transportClient.TLSClientConfig.InsecureSkipVerify = true
}
