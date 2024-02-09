package core

import (
	"context"
	"net"
	"net/http"
)

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
