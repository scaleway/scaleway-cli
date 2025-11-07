//go:build !wasm

package gotty

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/containerd/console"
	"github.com/gorilla/websocket"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	// GoTTY ingress message code
	outputCode         = '1'
	setWindowTitleCode = '3'

	// GoTTY egress message code
	inputCode          = '1'
	pingCode           = '2'
	resizeTerminalCode = '3'
)

type Client struct {
	wsURL     string
	serverID  string
	secretKey string
}

// NewClient returns a GoTTY client.
func NewClient(zone scw.Zone, serverID string, secretKey string) (*Client, error) {
	return &Client{
		wsURL:     fmt.Sprintf("wss://tty.%s.scaleway.com/v2/ws", zone.String()),
		serverID:  serverID,
		secretKey: secretKey,
	}, nil
}

func (c *Client) Connect() error {
	wsDialer := websocket.Dialer{}
	conn, _, err := wsDialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}
	defer func() {
		// Websocket protocol require the server to close the connection.
		// This sent a close request.
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
	}()

	// This is how scaleway implement gotty authentication
	err = conn.WriteJSON(map[string]string{
		"AuthToken": "",
		"Arguments": "?" + url.Values{"arg": []string{c.secretKey, c.serverID}}.Encode(),
	})
	if err != nil {
		return fmt.Errorf("failed to send auth json: %w", err)
	}

	cns, err := console.ConsoleFromFile(os.Stdin)
	if err != nil {
		return fmt.Errorf("os.Stdin doesn't seems to be a valid terminal: %w", err)
	}
	err = cns.SetRaw()
	if err != nil {
		return fmt.Errorf("error setting raw terminal: %w", err)
	}
	defer cns.Reset() //nolint:errcheck
	defer cns.Close()

	// Create a channel that will receive all resizes signals
	resizeChan := make(chan bool, 1)
	unsubscribe := subscribeToResize(resizeChan)
	defer unsubscribe()

	wsChan, wsErrChan := websocketReader(conn)
	cnsChan, cnsErrChan := consoleReader(cns)

	// Force first resize
	resizeChan <- true

	for {
		select {
		// Resize event: we send new terminal size to the server
		case <-resizeChan:
			size, err := cns.Size()
			if err != nil {
				return err
			}
			message := fmt.Sprintf(
				`%c{"columns":%d,"rows":%d}`,
				resizeTerminalCode,
				size.Width,
				size.Height,
			)
			err = conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				return fmt.Errorf("failed to write message on websocket: %w", err)
			}

		// We receive a message from the server
		case message := <-wsChan:
			// If message is empty, connection was closed
			if len(message) == 0 {
				return nil
			}

			switch message[0] {
			// The message contain data that should be printed
			case outputCode:
				buf, err := base64.StdEncoding.DecodeString(string(message[1:]))
				if err != nil {
					return fmt.Errorf("failed to decode base64 output payload: %w", err)
				}
				os.Stdout.Write(buf)
			// The message contain a new terminal title
			case setWindowTitleCode:
				newTitle := string(message[1:])
				fmt.Fprintf(os.Stdout, "\033]0;%s\007", newTitle)
			// We ignore other type of events
			default:
			}

		// We read something on the console (probably user input): we send it to the server.
		case message := <-cnsChan:
			// If message is empty the console has been closed
			if len(message) == 0 {
				return nil
			}
			err = conn.WriteMessage(websocket.TextMessage, append([]byte{inputCode}, message...))
			if err != nil {
				return fmt.Errorf("failed to write message on websocket: %w", err)
			}

		// We make sure to send a ping every 30s to keep the connection alive.
		case <-time.After(30 * time.Second):
			err = conn.WriteMessage(websocket.TextMessage, []byte{pingCode})
			if err != nil {
				return fmt.Errorf("failed to ping websocket: %w", err)
			}

		// If we receive an error from one of the 2 reader we return it
		case err := <-wsErrChan:
			if err != nil {
				return fmt.Errorf("websocket reader error: %w", err)
			}

			return nil
		case err := <-cnsErrChan:
			if err != nil {
				return fmt.Errorf("console reader error: %w", err)
			}

			return nil
		}
	}
}

// websocketReader start a go routine to read incoming messages on the websocket.
// It return 2 channels one with read message and one in case of errors.
func websocketReader(conn *websocket.Conn) (chan []byte, chan error) {
	readChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		defer close(readChan)
		defer close(errChan)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				// if a error is a close error exit properly
				if _, isClose := err.(*websocket.CloseError); !isClose {
					errChan <- err
				}

				return
			}
			readChan <- data
		}
	}()

	return readChan, errChan
}

// consoleReader start a go routine to read user input on the console.
// It return 2 channels one with read input and one in case of errors.
func consoleReader(cns console.Console) (chan []byte, chan error) {
	readChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		defer close(readChan)
		defer close(errChan)

		for {
			buff := make([]byte, 128)
			size, err := cns.Read(buff)
			if err != nil {
				errChan <- err

				return
			}
			readChan <- buff[:size]
		}
	}()

	return readChan, errChan
}
