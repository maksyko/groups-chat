package models

import (
	"github.com/twinj/uuid"
	"strings"

	"github.com/pkg/errors"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/datastore"
	"gopkg.in/mgo.v2/bson"
)

type usersHandler struct {
	UsernameRE string
}

func NewUsersHandler() *usersHandler {
	return &usersHandler{
		UsernameRE: "^[A-Za-z0-9_]{1,15}$", // username length copied from Twitter
	}
}

func (h *usersHandler) Register(username, email, password string) (*User, error) {
	u := &User{
		ID:        uuid.NewV4().String(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: app.Timestamp(),
		UpdatedAt: app.Timestamp(),
	}

	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if err := h.UsernameValid(u); err != nil {
		return nil, err
	}

	if err := h.EmailValid(u); err != nil {
		return nil, err
	}

	if err := h.PasswordValid(u); err != nil {
		return nil, err
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hpass)

	if err := u.Create(); err != nil {
		return nil, err
	}

	return u, nil
}

func (h *usersHandler) UsernameValid(u *User) error {
	if u.Username == "" {
		return errors.New("username required")
	}

	//if matched, _ := regexp.MatchString(h.UsernameRE, u.Username); !matched {
	//	return errors.New("username invalid")
	//}

	if u.UsernameAvailable() {
		return errors.New("username exists")
	}
	return nil
}

func (h *usersHandler) EmailValid(u *User) error {
	if u.Email == "" {
		return errors.New("email required")
	}

	if !govalidator.IsEmail(u.Email) {
		return errors.New("email invalid")
	}

	if u.EmailAvailable() {
		return errors.New("email exists")
	}
	return nil
}

func (h *usersHandler) PasswordValid(u *User) error {
	if u.Password == "" {
		return errors.New("password required")
	}
	return nil
}

func (h *usersHandler) ByUsername(username string) (*User, error) {
	var u *User
	return u, datastore.DB.Users.Find(bson.M{"username": strings.ToLower(username)}).One(&u)
}

func (h *usersHandler) ByEmail(email string) (*User, error) {
	var u *User
	return u, datastore.DB.Users.Find(bson.M{"username": strings.ToLower(email)}).One(&u)
}

func (h *usersHandler) ByID(ID string) (*User, error) {
	var u *User
	return u, datastore.DB.Users.FindId(ID).One(&u)
}

func (h *usersHandler) Auth(userPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
