package rpc

import (
	"encoding/json"

	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/messaging/client"
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ievgen-ma/groups-chat/models"
)

var eventsService = newEvents()
var messageService = newMessages()
var typingService = newTyping()

func CallMethod(c *client.Client, r *protocol.RPC) {
	processMsg := func(obj interface{}) {
		byteData, _ := json.Marshal(r.Body)
		err := json.Unmarshal(byteData, obj)
		if err != nil {
			app.Logger.Errorf("MESSAGING: Unable to read message %v %v\n", r.Body, err)
		}
	}

	switch r.Method {
	case protocol.RPC_MESSAGE_GET:
		params := protocol.RpcMessageGet{}
		processMsg(&params)
		eventsService.Get(c, &params)

	case protocol.RPC_MESSAGE_SEND:
		params := protocol.RpcMessageSend{}
		processMsg(&params)
		messageService.Send(c, &params)

	case protocol.RPC_MESSAGE_DELIVERED:
		params := protocol.RpcMessageDelivered{}
		processMsg(&params)
		messageService.Delivered(c, &params)

	case protocol.RPC_MESSAGE_READ:
		params := protocol.RpcMessageRead{}
		processMsg(&params)
		messageService.Read(c, &params)

	case protocol.RPC_TYPING_START:
		params := protocol.RpcTypingStart{}
		processMsg(&params)
		typingService.Start(c, &params)

	case protocol.RPC_TYPING_END:
		params := protocol.RpcTypingEnd{}
		processMsg(&params)
		typingService.End(c, &params)
	}
}

func withGroup(groupID, userID string, f func(group *models.Group)) {
	g, err := models.Groups.ByIDAndUserID(groupID, userID)
	if err != nil {
		return
	}
	f(g)
}
