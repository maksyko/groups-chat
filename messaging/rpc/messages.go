package rpc

import (
	"github.com/ievgen-ma/groups-chat/messaging/client"
	"github.com/ievgen-ma/groups-chat/messaging/events"
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/models"
)

type messages struct{}

func newMessages() *messages {
	return &messages{}
}

func (m *messages) Send(c *client.Client, p *protocol.RpcMessageSend) {
	withGroup(p.GroupID, c.UserID, func(group *models.Group) {
		msg, _ := models.Messages.Create(p.GroupID, c.UserID, p.Data, app.Timestamp())

		e := events.NewMessage(msg)
		e.SaveForUsers(msg.ID, group.UserIDs)

		es := events.NewMessageSent(msg.ID, e.Timestamp)
		es.SaveForUser(msg.ID, group.UserID)
		es.SendToUser(group.UserID)

		e.SendToUsersWithoutMe(c.SessionID, group.UserIDs)
	})
}

func (m *messages) Delivered(c *client.Client, p *protocol.RpcMessageDelivered) {
	msg, err := models.Messages.ByID(p.MessageID)
	if err != nil {
		return
	}

	withGroup(msg.GroupID, c.UserID, func(group *models.Group) {
		e := events.NewMessageDelivered(msg.ID, app.Timestamp())
		e.SaveForUser(msg.ID, msg.UserID)
		e.SendToUsersWithoutMe(c.SessionID, []string{msg.UserID, c.UserID})
		e.DeleteOldEvents(msg.ID)
	})
}

func (m *messages) Read(c *client.Client, p *protocol.RpcMessageRead) {
	msg, err := models.Messages.ByID(p.MessageID)
	if err != nil {
		return
	}

	withGroup(msg.GroupID, c.UserID, func(group *models.Group) {
		e := events.NewMessageRead(msg.ID, app.Timestamp())
		e.SaveForUser(msg.ID, msg.UserID)
		e.SendToUsersWithoutMe(c.SessionID, []string{msg.UserID, c.UserID})
		e.DeleteOldEvents(msg.ID)
	})
}
