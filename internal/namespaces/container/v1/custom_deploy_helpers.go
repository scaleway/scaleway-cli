//go:build darwin || linux || windows

package container

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	pack "github.com/buildpacks/pack/pkg/client"
	docker "github.com/moby/moby/client"
)

type DockerClient interface {
	pack.DockerClient

	ImagePush(
		ctx context.Context,
		image string,
		options docker.ImagePushOptions,
	) (docker.ImagePushResponse, error)
}

type CustomDockerClient struct {
	*docker.Client

	httpClient *http.Client
}

func NewCustomDockerClient(httpClient *http.Client) (*CustomDockerClient, error) {
	dockerClient, err := docker.New(
		docker.FromEnv,
		docker.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %w", err)
	}

	return &CustomDockerClient{
		Client:     dockerClient,
		httpClient: httpClient,
	}, nil
}

func (c *CustomDockerClient) ContainerAttach(
	_ context.Context,
	container string,
	options docker.ContainerAttachOptions,
) (docker.ContainerAttachResult, error) {
	query := url.Values{}
	if options.Stream {
		query.Set("stream", "1")
	}
	if options.Stdin {
		query.Set("stdin", "1")
	}
	if options.Stdout {
		query.Set("stdout", "1")
	}
	if options.Stderr {
		query.Set("stderr", "1")
	}
	if options.DetachKeys != "" {
		query.Set("detachKeys", options.DetachKeys)
	}
	if options.Logs {
		query.Set("logs", "1")
	}

	requestURL := &url.URL{
		Scheme:   "http",
		Host:     strings.TrimPrefix(c.DaemonHost(), "unix://"),
		Path:     fmt.Sprintf("/containers/%s/attach", container),
		RawQuery: query.Encode(),
	}

	reader, writer := net.Pipe()

	go func() {
		defer writer.Close()

		resp, err := c.httpClient.Do(&http.Request{
			Method:     http.MethodPost,
			Host:       "docker",
			URL:        requestURL,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: map[string][]string{
				"Content-Type": {"text/plain"},
				"Connection":   {"Upgrade"},
				"Upgrade":      {"tcp"},
			},
		})
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusSwitchingProtocols {
			panic(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		}

		_, err = io.Copy(writer, resp.Body)
		if err != nil {
			panic(err)
		}
	}()

	return docker.ContainerAttachResult{
		HijackedResponse: docker.NewHijackedResponse(reader, "text/plain"),
	}, nil
}
