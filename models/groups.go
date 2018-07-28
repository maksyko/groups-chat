package models

import (
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/twinj/uuid"
	"github.com/ievgen-ma/groups-chat/datastore"
	"gopkg.in/mgo.v2/bson"
)

type groups struct{}

var Groups = new(groups)

type Group struct {
	ID        string   `bson:"_id"`
	Name      string   `bson:"name"`
	UserID    string   `bson:"user_id"`
	UserIDs   []string `bson:"user_ids"`
	Archived  []string `bson:"archived"`
	Deleted   []string `bson:"deleted"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
}

func (groups) Create(name, userID string, userIDs []string) (*Group, error) {
	c := &Group{
		ID:        uuid.NewV4().String(),
		Name:      name,
		UserID:    userID,
		UserIDs:   userIDs,
		CreatedAt: app.Timestamp(),
		UpdatedAt: app.Timestamp(),
	}

	return c, datastore.DB.Groups.Insert(c)
}

func (groups) ByIDAndUserID(ID, userID string) (*Group, error) {
	var g *Group
	return g, datastore.DB.Groups.Find(bson.M{
		"_id":      ID,
		"user_ids": userID,
	}).One(&g)
}

func (groups) ByID(ID string) (*Group, error) {
	var g *Group
	return g, datastore.DB.Groups.FindId(ID).One(&g)
}

func (g *Group) UpdateUserIDs() error {
	return datastore.DB.Groups.UpdateId(g.ID, bson.M{
		"$set": bson.M{
			"user_ids":   g.UserIDs,
			"updated_at": app.Timestamp(),
		},
	})
}
