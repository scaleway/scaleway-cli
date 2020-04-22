package console

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

// GetAuthTokenURL transforms a GoTTY http URL to its AuthToken file URL
func GetAuthTokenURL(httpURL string) (*url.URL, *http.Header, error) {
	header := http.Header{}
	target, err := url.Parse(httpURL)
	if err != nil {
		return nil, nil, err
	}

	target.Path = strings.TrimLeft(target.Path+"auth_token.js", "/")

	if target.User != nil {
		header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(target.User.String())))
		target.User = nil
	}

	return target, &header, nil
}

// GetURLQuery returns url.query
func GetURLQuery(rawurl string) (url.Values, error) {
	target, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return target.Query(), nil
}

// GetWebsocketURL transforms a GoTTY http URL to its WebSocket URL
func GetWebsocketURL(httpURL string) (*url.URL, *http.Header, error) {
	header := http.Header{}
	target, err := url.Parse(httpURL)
	if err != nil {
		return nil, nil, err
	}

	if target.Scheme == "https" {
		target.Scheme = "wss"
	} else {
		target.Scheme = "ws"
	}

	target.Path = strings.TrimLeft(target.Path+"ws", "/")

	if target.User != nil {
		header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(target.User.String())))
		target.User = nil
	}

	return target, &header, nil
}
