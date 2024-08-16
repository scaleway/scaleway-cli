package webcallback

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/logger"
)

// WebCallback is a web server that will wait for a callback
type WebCallback struct {
	port int

	tokenChan chan string
	errChan   chan error
	srv       *http.Server
	listener  net.Listener
}

func New(opts ...Options) *WebCallback {
	wb := new(WebCallback)
	for _, opt := range opts {
		opt(wb)
	}

	return wb
}

func (wb *WebCallback) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wb.tokenChan = make(chan string)
	wb.errChan = make(chan error)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(wb.port))
	if err != nil {
		cancel()
		return err
	}
	wb.listener = listener
	wb.port = listener.Addr().(*net.TCPAddr).Port
	wb.srv = &http.Server{Addr: ":" + strconv.Itoa(wb.port)}

	wb.srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "callback") {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(webpageString("Invalid URL")))
		}
		token := r.URL.Query().Get("token")
		if token != "" {
			wb.tokenChan <- r.URL.Query().Get("token")
			w.WriteHeader(200)
			_, _ = w.Write([]byte(webpageString("You can close this page.")))
		} else {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(webpageString("Invalid Token.")))
		}

		cancel()
	})

	go func() {
		err = wb.srv.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			wb.errChan <- err
		}
		cancel()
	}()

	return nil
}

func (wb *WebCallback) Wait(ctx context.Context) (string, error) {
	defer wb.Close()
	select {
	case err := <-wb.errChan:
		return "", err
	case token := <-wb.tokenChan:
		return token, nil
	case <-ctx.Done():
		logger.Warningf("context canceled, closing web server")
		return "", ctx.Err()
	}
}

func (wb *WebCallback) Close() {
	err := wb.srv.Close()
	if err != nil {
		logger.Warningf("failed to close web server: %v", err)
	}
}

// Port returns the port used by the web server. It may be chosen randomly if let as default when starting server.
func (wb *WebCallback) Port() int {
	return wb.port
}

func webpageString(msg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
%s
</body>
</html>
`, msg)
}
