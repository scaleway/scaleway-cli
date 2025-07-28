package core

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/scaleway/scaleway-sdk-go/logger"
)

const defaultRetryInterval = 1 * time.Second

var socketTransport = &http.Transport{}

func init() {
	socketTransport.DisableCompression = true
	socketTransport.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
		return net.DialTimeout("unix", "/var/run/docker.sock", 32000000000)
	}
}

type SocketPassthroughTransport struct{}

func (r *SocketPassthroughTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.URL.Host == "/var/run/docker.sock" {
		return socketTransport.RoundTrip(request)
	}

	return http.DefaultTransport.RoundTrip(request)
}

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
		logger.Warningf(
			"cli: cannot use insecure mode with DefaultTransport of type %T",
			http.DefaultTransport,
		)

		return
	}
	if transportClient.TLSClientConfig == nil {
		transportClient.TLSClientConfig = &tls.Config{}
	}
	transportClient.TLSClientConfig.InsecureSkipVerify = true
}
