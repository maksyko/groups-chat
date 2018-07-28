package events

import (
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ievgen-ma/groups-chat/models"
)

func NewMessage(m *models.Message) *Event {
	event := NewEvent(protocol.EVENT_MESSAGE, m.CreatedAt)
	event.Body = protocol.EventMessage{
		MessageID: m.ID,
		Data:      m.Body.Data,
	}
	return event
}

func NewMessageSent(messageID string, ts int64) *Event {
	event := NewEvent(protocol.EVENT_MESSAGE_SENT, ts)
	event.Body = protocol.EventMessageSent{
		MessageID: messageID,
	}
	return event
}

func NewMessageDelivered(messageID string, ts int64) *Event {
	event := NewEvent(protocol.EVENT_MESSAGE_DELIVERED, ts)
	event.Body = protocol.EventMessageDelivered{
		MessageID: messageID,
	}
	return event
}

func NewMessageRead(messageID string, ts int64) *Event {
	event := NewEvent(protocol.EVENT_MESSAGE_READ, ts)
	event.Body = protocol.EventMessageSent{
		MessageID: messageID,
	}
	return event
}
