package rpc

import (
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ievgen-ma/groups-chat/messaging/client"
	"github.com/ievgen-ma/groups-chat/models"
	"github.com/ievgen-ma/groups-chat/messaging/events"
	"errors"
	"github.com/ievgen-ma/groups-chat/app"
)

type Events struct{}

func newEvents() *Events {
	return &Events{}
}

func (e *Events) Get(c *client.Client, p *protocol.RpcMessageGet) {
	es, err := models.Events.ByUserIDAndTimestamp(c.UserID, p.Timestamp)
	if err != nil {
		return
	}

	for _, event := range es {
		ev, err := e.execute(event)
		if err != nil {
			continue
		}

		if ev == nil {
			app.Logger.Infof("MESSAGING: Could not execute event %s with ID %s and message ID %s", event.Type, event.ID, event.ObjectID)
			continue
		}

		ev.SendToUser(c.UserID)
	}
}

func (e *Events) execute(event *models.Event) (*events.Event, error) {
	switch event.Type {
	case protocol.EVENT_MESSAGE:
		return e.executeMessage(event.ObjectID)
	case protocol.EVENT_MESSAGE_SENT:
		return e.executeMessageSent(event.ObjectID, event.Timestamp)
	case protocol.EVENT_MESSAGE_DELIVERED:
		return e.executeMessageDelivered(event.ObjectID, event.Timestamp)
	case protocol.EVENT_MESSAGE_READ:
		return e.executeMessageRead(event.ObjectID, event.Timestamp)
	case protocol.EVENT_GROUP:
		return e.executeGroup(event.ObjectID)
	case protocol.EVENT_GROUP_JOINED:
		return e.executeGroupJoined(event.ObjectID)
	case protocol.EVENT_GROUP_LEFT:
		return e.executeGroupLeft(event.ObjectID)
	}

	return nil, errors.New("wrong event type")
}

func (e *Events) executeMessage(messageID string) (*events.Event, error) {
	m, err := models.Messages.ByID(messageID)
	if err != nil {
		return nil, err
	}
	return events.NewMessage(m), nil
}

func (e *Events) executeMessageSent(messageID string, ts int64) (*events.Event, error) {
	m, err := models.Messages.ByID(messageID)
	if err != nil {
		return nil, err
	}
	return events.NewMessageSent(m.ID, ts), nil
}

func (e *Events) executeMessageDelivered(messageID string, ts int64) (*events.Event, error) {
	m, err := models.Messages.ByID(messageID)
	if err != nil {
		return nil, err
	}
	return events.NewMessageDelivered(m.ID, ts), nil
}

func (e *Events) executeMessageRead(messageID string, ts int64) (*events.Event, error) {
	m, err := models.Messages.ByID(messageID)
	if err != nil {
		return nil, err
	}
	return events.NewMessageRead(m.ID, ts), nil
}

func (e *Events) executeGroup(groupID string) (*events.Event, error) {
	g, err := models.Groups.ByID(groupID)
	if err != nil {
		return nil, err
	}
	return events.NewGroup(g), nil
}

func (e *Events) executeGroupJoined(groupID string) (*events.Event, error) {
	g, err := models.Groups.ByID(groupID)
	if err != nil {
		return nil, err
	}
	return events.NewGroupJoined(g), nil
}

func (e *Events) executeGroupLeft(groupID string) (*events.Event, error) {
	g, err := models.Groups.ByID(groupID)
	if err != nil {
		return nil, err
	}
	return events.NewGroupLeft(g), nil
}
