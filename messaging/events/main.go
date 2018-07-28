package events

import (
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/ievgen-ma/groups-chat/models"
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/messaging/client"
)

type Event struct {
	*protocol.Event
}

func NewEvent(t protocol.EventType, ts int64) *Event {
	pe := &protocol.Event{
		Type:      t,
		Timestamp: ts,
	}

	return &Event{pe}
}

func (e *Event) SendToUser(userID string) {
	e.SendToUsersWithoutMe("", []string{userID})
}

func (e *Event) SendToUsers(userIDs []string) {
	e.SendToUsersWithoutMe("", userIDs)
}

func (e *Event) SendToUsersWithoutMe(sessionID string, userIDs []string) {
	for _, uID := range userIDs {
		sessions, err := models.Sessions.ByUserID(uID)
		if err != nil {
			continue
		}
		for _, s := range sessions {
			if sessionID != s.ID {
				app.Logger.Debugf("MESSAGING: Send event to session %v: %v", s, e)
				c := client.NewFromSession(uID, s.ID)
				c.SendEvent(e.Event)
			}
		}
	}
}

func (e *Event) SaveForUser(objectID, userID string) {
	e.SaveForUsers(objectID, []string{userID})
}

func (e *Event) SaveForUsers(objectID string, userIDs []string) {
	models.Events.Create(e.Type, objectID, userIDs, e.Timestamp)
}

func (e *Event) DeleteOldEvents(objectID string) {
	models.Events.DeleteOldEvents(objectID, e.Type, e.Timestamp)
}
