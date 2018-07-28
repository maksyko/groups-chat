package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"github.com/ievgen-ma/groups-chat/models"
)

type SessionsController struct {
	BaseController
}

func NewSessionsController() *SessionsController {
	return &SessionsController{}
}

var sessions = models.NewSessionsHandler()

func (c *SessionsController) Create(w rest.ResponseWriter, r *rest.Request) {
	if err := c.Authenticate(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	in := struct {
		DeviceID string `json:"device_id"`
		Platform string `json:"platform"`
		Model    string `json:"model"`
		Build    int    `json:"build"`
		Name     string `json:"name"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	s, err := sessions.Create(c.User.ID, in.DeviceID, in.Platform, in.Model, in.Build, in.Name, r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(map[string]interface{}{
		"id":            s.ID,
		"messaging_url": s.MessagingURL(),
		"created_at":    s.CreatedAt,
	})
}
