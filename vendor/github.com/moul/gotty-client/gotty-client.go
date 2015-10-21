package gottyclient

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/scaleway/scaleway-cli/vendor/github.com/moul/gotty-client/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/moul/gotty-client/vendor/github.com/gorilla/websocket"
	"github.com/scaleway/scaleway-cli/vendor/github.com/moul/gotty-client/vendor/golang.org/x/crypto/ssh/terminal"
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

type Client struct {
	Dialer     *websocket.Dialer
	Conn       *websocket.Conn
	URL        string
	Connected  bool
	WriteMutex *sync.Mutex
}

type querySingleType struct {
	AuthToken string `json:"AuthToken"`
	Arguments string `json:"Arguments"`
}

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

	logrus.Debugf("Fetching auth token auth-token: %q", target.String())
	req, err := http.NewRequest("GET", target.String(), nil)
	req.Header = *header
	client := http.Client{}
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
	if len(output) == 0 {
		return "", fmt.Errorf("Cannot fetch GoTTY auth-token, please upgrade your GoTTY server.")
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
	logrus.Debugf("Auth-token: %q", authToken)

	// Open WebSocket connection
	target, header, err := GetWebsocketURL(c.URL)
	if err != nil {
		return err
	}
	logrus.Debugf("Connecting to websocket: %q", target.String())
	conn, _, err := c.Dialer.Dial(target.String(), *header)
	if err != nil {
		return err
	}
	c.Conn = conn

	// Pass arguments and auth-token
	query, err := GetURLQuery(c.URL)
	if err != nil {
		return err
	}
	var querySingle querySingleType = querySingleType{
		Arguments: "?" + query.Encode(),
		AuthToken: authToken,
	}
	json, err := json.Marshal(querySingle)
	if err != nil {
		logrus.Errorf("Failed to parse init message %v", err)
		return err
	}
	// Send Json
	logrus.Debugf("Sending arguments and auth-token")
	err = c.write(json)
	if err != nil {
		return err
	}

	go c.pingLoop()

	return nil
}

func (c *Client) pingLoop() {
	for {
		logrus.Debugf("Sending ping")
		c.write([]byte("1"))
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

	done := make(chan bool)
	go c.readLoop(done)
	go c.writeLoop(done)
	go c.termsizeLoop(done)
	<-done
	return nil
}

type winsize struct {
	Rows    uint16 `json:"rows"`
	Columns uint16 `json:"columns"`
	// unused
	x uint16
	y uint16
}

func (c *Client) termsizeLoop(done chan bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	ws := winsize{}

	for {
		syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(0), uintptr(syscall.TIOCGWINSZ),
			uintptr(unsafe.Pointer(&ws)))

		b, err := json.Marshal(ws)
		if err != nil {
			logrus.Warnf("json.Marshal error: %v", err)
		}

		err = c.write(append([]byte("2"), b...))
		if err != nil {
			logrus.Warnf("ws.WriteMessage failed: %v", err)
		}

		<-ch
	}
}

func (c *Client) writeLoop(done chan bool) {
	oldState, err := terminal.MakeRaw(0)
	if err == nil {
		defer terminal.Restore(0, oldState)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		x, size, err := reader.ReadRune()
		if size <= 0 || err != nil {
			done <- true
			return
		}

		p := make([]byte, size)
		utf8.EncodeRune(p, x)

		err = c.write(append([]byte("0"), p...))
		if err != nil {
			done <- true
			return
		}
	}
}

func (c *Client) readLoop(done chan bool) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			done <- true
			logrus.Warnf("c.Conn.ReadMessage: %v", err)
			return
		}
		if len(data) == 0 {
			done <- true
			logrus.Warnf("An error has occured")
			return
		}
		switch data[0] {
		case '0': // data
			buf, err := base64.StdEncoding.DecodeString(string(data[1:]))
			if err != nil {
				logrus.Warnf("Invalid base64 content: %q", data[1:])
			}
			fmt.Print(string(buf))
		case '1': // pong
		case '2': // new title
			newTitle := string(data[1:])
			fmt.Printf("\033]0;%s\007", newTitle)
		case '3': // json prefs
			logrus.Debugf("Unhandled protocol message: json pref: %s", string(data[1:]))
		case '4': // autoreconnect
			logrus.Debugf("Unhandled protocol message: autoreconnect: %s", string(data))
		default:
			logrus.Warnf("Unhandled protocol message: %s", string(data))
		}
	}
}

// NewClient returns a GoTTY client object
func NewClient(httpURL string) (*Client, error) {
	return &Client{
		Dialer:     &websocket.Dialer{},
		URL:        httpURL,
		WriteMutex: &sync.Mutex{},
	}, nil
}
