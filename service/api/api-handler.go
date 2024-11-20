package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/test-users", rt.TestUsers)

	// login
	rt.router.POST("/session", rt.session)

	// find user
	rt.router.GET("/users/search", rt.findUser)

	return rt.router
}
