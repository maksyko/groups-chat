package client

import (
	"github.com/gorilla/websocket"
	"sync"
	"fmt"
)

type clientRegistry struct {
	mutex       sync.RWMutex
	connections map[string]*websocket.Conn
}

func newRegistry() *clientRegistry {
	return &clientRegistry{connections: map[string]*websocket.Conn{}}
}

func (c *clientRegistry) set(client *Client) {
	// set online

	c.mutex.Lock()
	c.connections[client.UserID] = client.Connection
	c.mutex.Unlock()

	// cache client connect
}

func (c *clientRegistry) get(ID string) (*websocket.Conn, error) {

	c.mutex.RLock()
	conn, ok := c.connections[ID]
	c.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("could not find client with client UserID %s", ID)
	}
	return conn, nil
}

func (c *clientRegistry) delete(client *Client) {
	// set offline

	c.mutex.Lock()
	delete(c.connections, client.UserID)
	c.mutex.Unlock()

	// cache client disconnect
}