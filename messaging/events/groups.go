package events

import (
	"github.com/ievgen-ma/groups-chat/models"
	"github.com/ievgen-ma/groups-chat/protocol"
)

func NewGroup(g *models.Group) *Event {
	event := NewEvent(protocol.EVENT_GROUP, g.CreatedAt)
	event.Body = protocol.EventGroup{
		GroupID: g.ID,
		Name:    g.Name,
		UserIDs: g.UserIDs,
	}
	return event
}

func NewGroupJoined(g *models.Group) *Event {
	event := NewEvent(protocol.EVENT_GROUP_JOINED, g.CreatedAt)
	event.Body = protocol.EventGroupJoined{
		GroupID: g.ID,
		UserID:  g.UserID,
	}
	return event
}

func NewGroupLeft(g *models.Group) *Event {
	event := NewEvent(protocol.EVENT_GROUP_LEFT, g.CreatedAt)
	event.Body = protocol.EventGroupLeft{
		GroupID: g.ID,
		UserID:  g.UserID,
	}
	return event
}
