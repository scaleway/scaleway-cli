package console

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

// Useful message types for gotty
const (
	Output         = '0'
	Pong           = '1'
	SetWindowTitle = '2'
	SetPreferences = '3'
	SetReconnect   = '4'

	Input          = '0'
	Ping           = '1'
	ResizeTerminal = '2'
)

type messageType struct {
	output         byte
	pong           byte
	setWindowTitle byte
	setPreferences byte
	setReconnect   byte
	input          byte
	ping           byte
	resizeTerminal byte
}

type Client struct {
	Dialer     *websocket.Dialer
	Conn       *websocket.Conn
	URL        string
	WriteMutex *sync.Mutex
	Output     io.Writer
	result     chan error
	Connected  bool
	EscapeKeys []byte
	message    *messageType
	WSOrigin   string
}

type querySingleType struct {
	AuthToken string `json:"AuthToken"`
	Arguments string `json:"Arguments"`
}

type winsize struct {
	Rows    uint16 `json:"rows"`
	Columns uint16 `json:"columns"`
	// unused
	x uint16
	y uint16
}
