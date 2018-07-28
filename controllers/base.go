package controllers

import (
	"github.com/ievgen-ma/groups-chat/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ievgen-ma/groups-chat/app/jwt"
)

type BaseController struct {
	User *models.User
}

var users = models.NewUsersHandler()

func (c *BaseController) Login(username, email, password string) (*models.User, error) {
	var (
		u   *models.User
		err error
	)

	if username != "" {
		u, err = users.ByUsername(username)
		if err != nil {
			return nil, err
		}
	} else if email != "" {
		u, err = users.ByEmail(email)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("username or email is required")
	}

	if !users.Auth(u.Password, password) {
		return nil, errors.New("password wrong")
	}

	return u, nil
}

func (c *BaseController) Authenticate(r *rest.Request) (err error) {
	token, err := jwt.Parse(r)
	if err != nil {
		return
	}

	u, err := users.ByID(token.UserID)
	if err != nil {
		return
	}

	c.User = u
	return nil
}
