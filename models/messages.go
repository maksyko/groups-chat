package models

import (
	"github.com/ievgen-ma/groups-chat/datastore"
	"github.com/twinj/uuid"
)

type Message struct {
	ID        string `bson:"_id"`
	UserID    string `bson:"user_id"`
	GroupID   string `bson:"group_id"`
	Body      Body   `bson:"body"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

type Body struct {
	Data string `bson:"data"`
}

type messages struct{}

var Messages = new(messages)

func (messages) ByID(ID string) (*Message, error) {
	var m *Message
	return m, datastore.DB.Messages.FindId(ID).One(&m)
}

func (messages) Create(groupID, userID, data string, ts int64) (*Message, error) {
	m := &Message{
		ID:        uuid.NewV4().String(),
		UserID:    userID,
		GroupID:   groupID,
		Body:      Body{Data: data},
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	return m, datastore.DB.Messages.Insert(m)
}
