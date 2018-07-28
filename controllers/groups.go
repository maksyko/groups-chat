package controllers

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"github.com/ievgen-ma/groups-chat/models"
	"github.com/ievgen-ma/groups-chat/messaging/events"
	"github.com/ievgen-ma/groups-chat/app"
)

type GroupsController struct {
	BaseController
}

func NewGroupController() *GroupsController {
	return &GroupsController{}
}

func (c *GroupsController) Create(w rest.ResponseWriter, r *rest.Request) {
	if err := c.Authenticate(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}
	in := struct {
		Name    string   `json:"name"`
		UserIDs []string `json:"user_ids"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	in.UserIDs = append(in.UserIDs, c.User.ID)

	g, err := models.Groups.Create(in.Name, c.User.ID, in.UserIDs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	if eg := events.NewGroup(g); eg != nil {
		eg.SaveForUsers(g.ID, g.UserIDs)
		eg.SendToUsers(g.UserIDs)
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(map[string]string{"id": g.ID})
}

func (c *GroupsController) Join(w rest.ResponseWriter, r *rest.Request) {
	if err := c.Authenticate(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}
	in := struct {
		UserID string `json:"user_id"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	g, err := models.Groups.ByID(r.PathParam("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	if app.SliceContains(g.UserIDs, in.UserID) {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": "user_id exists"})
		return
	}

	g.UserIDs = append(g.UserIDs, in.UserID)
	if err := g.UpdateUserIDs(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	if eg := events.NewGroupJoined(g); eg != nil {
		eg.SaveForUsers(g.ID, g.UserIDs)
		eg.SendToUsers(g.UserIDs)
	}
}

func (c *GroupsController) Left(w rest.ResponseWriter, r *rest.Request) {
	if err := c.Authenticate(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}
	in := struct {
		UserID string `json:"user_id"`
	}{}

	if err := r.DecodeJsonPayload(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	g, err := models.Groups.ByID(r.PathParam("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	if !app.SliceContains(g.UserIDs, in.UserID) {
		w.WriteHeader(http.StatusNotFound)
		w.WriteJson(map[string]string{"error": "user_id not found"})
		return
	}

	g.UserIDs = app.RemoveFromSlice(g.UserIDs, in.UserID)
	if err := g.UpdateUserIDs(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	if eg := events.NewGroupLeft(g); eg != nil {
		eg.SaveForUsers(g.ID, g.UserIDs)
		eg.SendToUsers(g.UserIDs)
	}
}
