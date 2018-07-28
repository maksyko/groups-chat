package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/ievgen-ma/groups-chat/datastore"
	"github.com/ievgen-ma/groups-chat/protocol"
	"github.com/twinj/uuid"
)

type Event struct {
	ID        string             `bson:"_id"`
	Type      protocol.EventType `bson:"type"`
	ObjectID  string             `bson:"object_id"`
	UserIDs   []string           `bson:"user_ids"`
	Timestamp int64              `bson:"timestamp"`
}

type events struct{}

var Events = new(events)

func (events) ByUserIDAndTimestamp(ID string, ts int64) ([]*Event, error) {
	var es []*Event
	return es, datastore.DB.Events.Find(bson.M{
		"user_ids":  ID,
		"timestamp": bson.M{"$gt": ts},
	}).Sort("timestamp").All(&es)
}

func (events) Create(typ protocol.EventType, objectID string, clientIDs []string, ts int64) (*Event, error) {
	e := &Event{
		ID:        uuid.NewV4().String(),
		Type:      typ,
		ObjectID:  objectID,
		UserIDs:   clientIDs,
		Timestamp: ts,
	}

	return e, datastore.DB.Events.Insert(e)
}

func (events) DeleteOldEvents(objectID string, typ protocol.EventType, ts int64) {
	datastore.DB.Events.RemoveAll(bson.M{
		"object_id": objectID,
		"timestamp": bson.M{"$lt": ts},
		"type":      typ,
	})
}
