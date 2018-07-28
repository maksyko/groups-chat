package datastore

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisHelper struct {
	Conn *redis.Pool
}

var Redis *redisHelper

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: time.Second * 180,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://127.0.0.1:6379/0")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	Redis = &redisHelper{Conn: newPool()}
}

func (r *redisHelper) Publish(channel string, data []byte) error {
	conn := r.Conn.Get()
	defer conn.Close()

	channel = fmt.Sprintf("messaging:%s", channel)
	return conn.Send("PUBLISH", channel, data)
}

func (r *redisHelper) Subscribe(f func(string, []byte)) error {
	conn := r.Conn.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}

	pattern := "messaging:*"
	err := psc.PSubscribe(pattern)
	defer psc.PUnsubscribe(pattern)
	defer psc.Close()
	if err != nil {
		return err
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.PMessage:
			f(v.Channel, v.Data)
		case error:
			return v
		}
	}
}
