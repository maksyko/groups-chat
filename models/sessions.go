package models

import (
	"github.com/ievgen-ma/groups-chat/datastore"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type sessions struct{}

var Sessions = new(sessions)

type Session struct {
	ID          string `bson:"_id"`
	UserID      string `bson:"user_id"`
	Type        string `bson:"type"`
	DeviceID    string `bson:"device_id"`
	Platform    string `bson:"platform"`
	Model       string `bson:"model"`
	IPAddress   string `bson:"ip_address"`
	Build       int    `bson:"build"`
	Name        string `bson:"name"`
	AccessToken string `bson:"access_token"`
	Online      bool   `bson:"online"`
	CreatedAt   int64  `bson:"created_at"`
	UpdatedAt   int64  `bson:"updated_at"`
	OnlineAt    int64  `bson:"online_at"`
	OfflineAt   int64  `bson:"offline_at"`
}

func (s *Session) Create() error {
	return datastore.DB.Sessions.Insert(s)
}

func (s *Session) MessagingURL() string {
	return fmt.Sprintf("%s/%s", "ws://127.0.0.1:3000", s.AccessToken)
}

func (sessions) ByUserID(userID string) ([]*Session, error) {
	var s []*Session
	return s, datastore.DB.Sessions.Find(bson.M{"user_id": userID}).All(&s)
}

func (sessions) ByAccessToken(accessToken string) (*Session, error) {
	var s *Session
	return s, datastore.DB.Sessions.Find(bson.M{"access_token": accessToken}).One(&s)
}
