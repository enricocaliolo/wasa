package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/test-users", rt.TestUsers)

	rt.router.POST("/session", rt.session)

	return rt.router
}
