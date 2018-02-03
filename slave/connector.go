package slave

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients map[string]*Client
}

// Client a client for communicating via websockets.
type Client struct {
	interrupt  chan struct{}
	done       chan struct{}
	socketConn *websocket.Conn
}

// NewClient creates a new client.
func NewClient(url url.URL) (*Client, error) {

	var err error
	ret := &Client{
		interrupt: make(chan struct{}),
		done:      make(chan struct{}),
	}

	ret.socketConn, _, err = websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}

	go func() {
		defer ret.socketConn.Close()
		defer close(ret.done)
		for {
			msgType, message, err := ret.socketConn.ReadMessage()
			if err != nil || msgType == websocket.CloseMessage {
				return
			}

			switch msgType {
			case websocket.CloseMessage:
				return
			case websocket.TextMessage:
				// todo text message handle
				break
			case websocket.PingMessage:
				ret.socketConn.WriteMessage(websocket.PongMessage, message)
			case websocket.PongMessage:
			case websocket.BinaryMessage:
			}

		}
	}()

	return nil, nil
}

// Close closes the websocket client
func (c *Client) Close() error {
	err := c.socketConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	select {
	case <-c.done:
	case <-time.After(time.Second * 2):
	}
	return c.Close()
}
