package main

import (
	"net/http"

	"github.com/ievgen-ma/groups-chat/app"
	"github.com/ievgen-ma/groups-chat/controllers"
	"github.com/ievgen-ma/groups-chat/messaging"
	"github.com/ant0ine/go-json-rest/rest"
)

var (
	users    = controllers.NewUsersController()
	groups   = controllers.NewGroupController()
	sessions = controllers.NewSessionsController()
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	r, err := rest.MakeRouter(
		// LOGIN RESOURCE
		rest.Post("/login", users.Login),
		// USERS RESOURCE
		rest.Post("/users", users.Register),
		// GROUPS RESOURCE
		rest.Post("/groups", groups.Create),
		rest.Put("/groups/:id/join", groups.Join),
		rest.Put("/groups/:id/left", groups.Left),
		// SESSIONS RESOURCE
		rest.Post("/sessions", sessions.Create),
		// MESSAGING RESOURCE
		rest.Get("/:access_token", messaging.Start),
	)
	if err != nil {
		app.Logger.Fatal(err)
	}
	api.SetApp(r)
	http.ListenAndServe(":3000", api.MakeHandler())
}
