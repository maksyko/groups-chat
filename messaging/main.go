package messaging

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/datastore"
	"github.com/ievgen-ma/groups-chat/messaging/client"
	"github.com/ievgen-ma/groups-chat/messaging/rpc"
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ievgen-ma/groups-chat/models"
)

const (
	writeWait  = 30 * time.Second
	pingPeriod = 10 * time.Second
)

func init() {
	go publishListener()
}

func publishListener() {
	datastore.Redis.Subscribe(func(channel string, data []byte) {
		chunks := strings.Split(channel, ":")
		sessionID := chunks[len(chunks)-1]
		conn, err := client.ConnectionBySessionID(sessionID)
		if err != nil {
			//app.Logger.Errorf("MESSAGING: Subscribe error %s %v %v", clientID, string(data), err)
			return
		}
		app.Logger.Infof("MESSAGING: Subscribe %s %v", sessionID, string(data))
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		conn.WriteMessage(websocket.TextMessage, data)
	})
}

func Start(w rest.ResponseWriter, r *rest.Request) {
	s, err := models.Sessions.ByAccessToken(r.PathParam("access_token"))
	if err != nil {
		return
	}

	c := client.NewFromRequest(s, w, r)
	if c == nil {
		app.Logger.Error("Unauthorized")
		return
	}

	ch := make(chan *protocol.RPC)
	go dispatcher(c, ch)
	reader(c, ch)
	close(ch)
}

func dispatcher(c *client.Client, ch chan *protocol.RPC) {
	app.Logger.Infof("MESSAGING: Dispatcher started for client ID %s", c.UserID)
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		app.Logger.Infof("MESSAGING: Disconnect client ID %s", c.UserID)
		ticker.Stop()
	}()

	for {
		select {
		case r, ok := <-ch:
			if !ok {
				app.Logger.Infof("MESSAGING: Could not receive event from client ID %s", c.UserID)
				c.SendCloseConnection()
				return
			}

			rpc.CallMethod(c, r)
		case <-ticker.C:
			err := c.SendPing()
			if err != nil {
				app.Logger.Infof("MESSAGING: Could not send ping to client ID %s %s", c.UserID, err)
				return
			}
		}
	}
}

func reader(c *client.Client, ch chan *protocol.RPC) {
	defer func() {
		c.Close()
		app.Logger.Infof("MESSAGING: Disconnect reader for client ID %s", c.UserID)
	}()

	c.Setup()

	for {
		rpc, err := c.ReadRPC()
		if err != nil {
			app.Logger.Infof("MESSAGING: Error in event from client ID %s: %v", c.UserID, err)
			return
		}

		app.Logger.Infof("MESSAGING: Received event from client ID %s: %v", c.UserID, rpc)

		ch <- rpc
	}
}
