package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	"github.com/ievgen-ma/groups-chat/app/jwt"
)

type UsersController struct {
	BaseController
}

func NewUsersController() *UsersController {
	return &UsersController{}
}

func (c *UsersController) Login(w rest.ResponseWriter, r *rest.Request) {
	in := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	u, err := c.BaseController.Login(in.Username, in.Email, in.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(map[string]string{"token": jwt.Create(u.ID)})
}
func (c *UsersController) Register(w rest.ResponseWriter, r *rest.Request) {
	in := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	u, err := users.Register(in.Username, in.Email, in.Password);
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(map[string]string{"id": u.ID})
}
