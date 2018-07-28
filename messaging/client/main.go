package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/datastore"
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ievgen-ma/groups-chat/models"
)

const (
	writeWait      = 30 * time.Second
	pongWait       = 30 * time.Second
	maxMessageSize = 1024 * 1024
)

var registry = newRegistry()

type Client struct {
	UserID     string
	SessionID  string
	Connection *websocket.Conn
}

func NewFromRequest(s *models.Session, w rest.ResponseWriter, r *rest.Request) *Client {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w.(http.ResponseWriter), r.Request, nil)
	if err != nil {
		return nil
	}

	client := &Client{
		UserID:     s.UserID,
		SessionID:  s.ID,
		Connection: conn,
	}

	registry.set(client)

	return client
}

func NewFromSession(userID, sessionID string) *Client {
	return &Client{
		UserID:    userID,
		SessionID: sessionID,
	}
}

func ConnectionBySessionID(ID string) (*websocket.Conn, error) {
	return registry.get(ID)
}

func (c *Client) SendCloseConnection() {
	c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
	c.Connection.WriteMessage(websocket.CloseMessage, nil)
}

func (c *Client) SendPing() error {
	c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Connection.WriteMessage(websocket.PingMessage, nil)
}

func (c *Client) Setup() {
	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		// set expire for client
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func (c *Client) Close() error {
	registry.delete(c)
	return c.Connection.Close()
}

func (c *Client) ReadRPC() (*protocol.RPC, error) {
	mt, data, err := c.Connection.ReadMessage()
	if err != nil {
		return nil, err
	}

	if mt == 0 {
		return nil, errors.New("invalid data received")
	}

	rpc := &protocol.RPC{}
	return rpc, json.Unmarshal(data, rpc)
}

func (c *Client) SendEvent(event *protocol.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		app.Logger.Error(err)
		return err
	}
	return datastore.Redis.Publish(c.UserID, data)
}
