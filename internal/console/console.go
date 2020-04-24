// Copied and modified from https://github.com/moul/gotty-client
// Copyright (c) 2015 Manfred Touron

package console

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/containerd/console"
	"github.com/gorilla/websocket"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

func (c *Client) write(data []byte) error {
	c.WriteMutex.Lock()
	defer c.WriteMutex.Unlock()
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

// GetAuthToken retrieves an Auth Token from dynamic auth_token.js file
func (c *Client) GetAuthToken() (string, error) {
	target, header, err := GetAuthTokenURL(c.URL)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", target.String(), nil)
	req.Header = *header
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case 200:
		// Everything is OK
	default:
		return "", fmt.Errorf("unknown status code: %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("var gotty_auth_token = '(.*)'")
	output := re.FindStringSubmatch(string(body))
	if len(output) == 0 || len(output) == 1 {
		return "", fmt.Errorf("cannot fetch console auth-token")
	}

	return output[1], nil
}

// Connect tries to dial a websocket server
func (c *Client) Connect() error {
	// Retrieve AuthToken
	authToken, err := c.GetAuthToken()
	if err != nil {
		return err
	}

	// Open WebSocket connection
	target, header, err := GetWebsocketURL(c.URL)
	if err != nil {
		return err
	}
	if c.WSOrigin != "" {
		header.Add("Origin", c.WSOrigin)
	}
	conn, _, err := c.Dialer.Dial(target.String(), *header)
	if err != nil {
		return err
	}
	c.Conn = conn
	c.Connected = true

	// Pass arguments and auth-token
	query, err := GetURLQuery(c.URL)
	if err != nil {
		return err
	}
	querySingle := querySingleType{
		Arguments: "?" + query.Encode(),
		AuthToken: authToken,
	}
	json, err := json.Marshal(querySingle)
	if err != nil {
		return err
	}
	// Send Json
	err = c.write(json)
	if err != nil {
		return err
	}

	// Initialize message types for gotty
	c.initMessageType()

	go c.pingLoop()

	return nil
}

// initMessageType initialize message types for gotty
func (c *Client) initMessageType() {
	c.message = &messageType{
		output:         Output,
		pong:           Pong,
		setWindowTitle: SetWindowTitle,
		setPreferences: SetPreferences,
		setReconnect:   SetReconnect,
		input:          Input,
		ping:           Ping,
		resizeTerminal: ResizeTerminal,
	}
}

func (c *Client) pingLoop() {
	for {
		c.write([]byte{c.message.ping})
		time.Sleep(30 * time.Second)
	}
}

// Close will nicely close the dialer
func (c *Client) Close() {
	c.Conn.Close()
}

// Loop will look indefinitely for new messages
func (c *Client) Loop() error {
	if !c.Connected {
		err := c.Connect()
		if err != nil {
			return err
		}
	}
	term, err := console.ConsoleFromFile(os.Stdout)
	if err != nil {
		return fmt.Errorf("os.Stdout is not a valid terminal")
	}
	err = term.SetRaw()
	if err != nil {
		return fmt.Errorf("error setting raw terminal: %v", err)
	}
	defer term.Reset()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go c.termsizeLoop(wg)

	wg.Add(1)
	go c.readLoop(wg)

	wg.Add(1)
	go c.writeLoop(wg)

	err = <-c.result
	if err != nil { // chan is closed only on user action
		close(c.result)
	}

	/* Wait for all of the above goroutines to finish */
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) termsizeLoop(wg *sync.WaitGroup) {
	defer wg.Done()

	ch := make(chan os.Signal, 1)
	notifySignalSIGWINCH(ch)
	defer resetSignalSIGWINCH()

	for {
		if b, err := syscallTIOCGWINSZ(); err != nil {
			logger.Errorf("failed to call TIOCGWINSZ: %s", err)
		} else {
			if err = c.write(append([]byte{c.message.resizeTerminal}, b...)); err != nil {
				logger.Errorf("failed to write: %s", err)
			}
		}
		select {
		case <-c.result:
			return
		case <-ch:
		}
	}
}

func (c *Client) writeLoop(wg *sync.WaitGroup) {
	logger.Debugf("starting writeLoop")
	defer wg.Done()

	// _ = c.write(append([]byte{c.message.input}, byte(0x0d))) TODO add?

	var size int
	buff := make([]byte, 128)

	reader := io.ReadCloser(os.Stdin)

	msgChan := make(chan error)

	defer reader.Close()

	for {
		go func() {
			var err error
			size, err = reader.Read(buff)
			msgChan <- err
		}()

		select {
		case <-c.result:
			return
		case err := <-msgChan:
			if err != nil {
				if err == io.EOF {
					err := c.write(append([]byte{c.message.input}, byte(4)))
					if err != nil {
						c.result <- fmt.Errorf("could not send EOF to console: %s", err)
						return
					}
					return
				}
				c.result <- fmt.Errorf("could not read from stdin: %s", err)
				return
			}
			if size <= 0 {
				continue
			}

			data := buff[:size]
			err = c.write(append([]byte{c.message.input}, data...))
			if err != nil {
				c.result <- fmt.Errorf("could not write to serial: %s", err)
				return
			}
		}
	}
}

func (c *Client) readLoop(wg *sync.WaitGroup) {
	defer wg.Done()

	msgChan := make(chan error)
	var data []byte

	for {
		go func() {
			var err error
			_, data, err = c.Conn.ReadMessage()
			msgChan <- err
		}()

		select {
		case <-c.result:
			return
		case err := <-msgChan:
			if err != nil {
				var closeError *websocket.CloseError
				if errors.As(err, &closeError) {
					if closeError.Code == websocket.CloseAbnormalClosure { // error received on Ctrl+Q
						close(c.result)
						return
					}
				}
				c.result <- fmt.Errorf("could not read from serial: %s", err)
				return
			}
			if len(data) == 0 {
				c.result <- fmt.Errorf("empty read from serial")
				return
			}
			switch data[0] {
			case c.message.output: // data
				buf, err := base64.StdEncoding.DecodeString(string(data[1:]))
				if err != nil {
					c.result <- fmt.Errorf("could not decode base64 serial output: %v", err)
					return
				}
				c.Output.Write(buf)
			case c.message.pong: // pong
			case c.message.setWindowTitle: // new title
				newTitle := string(data[1:])
				fmt.Fprintf(c.Output, "\033]0;%s\007", newTitle)
			case c.message.setPreferences: // json prefs
			case c.message.setReconnect: // autoreconnect
			default:
			}
		}
	}
}

// ParseURL parses an URL which may be incomplete and tries to standardize it
func ParseURL(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	switch parsed.Scheme {
	case "http", "https":
		// everything is ok
	default:
		return "", fmt.Errorf("missing scheme in URL")
	}
	return parsed.String(), nil
}

// NewClient returns a Console client object
func NewClient(consoleURL string, serverID string, secretKey string) (*Client, error) {
	url := fmt.Sprintf("%s/?arg=%s&arg=%s", consoleURL, secretKey, serverID)

	url, err := ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &Client{
		Dialer:     &websocket.Dialer{},
		URL:        url,
		WriteMutex: &sync.Mutex{},
		Output:     os.Stdout,
		result:     make(chan error),
	}, nil
}
