package datastore

import (
	"gopkg.in/mgo.v2"
)

type mongo struct {
	Events   *mgo.Collection
	Messages *mgo.Collection
	Groups   *mgo.Collection
	Users    *mgo.Collection
	Sessions *mgo.Collection
}

var DB *mongo

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := session.DB("ws-chat")

	DB = &mongo{
		Events:   db.C("events"),
		Messages: db.C("messages"),
		Groups:   db.C("groups"),
		Users:    db.C("users"),
		Sessions: db.C("sessions"),
	}
}
