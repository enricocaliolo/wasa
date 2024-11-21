package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *APIRouter) Handler() http.Handler {
	// user operations
	rt.router.GET("/test-users", rt.authMiddleware(rt.TestUsers))
	rt.router.PUT("/session", rt.login)
	rt.router.GET("/users/search", rt.authMiddleware(rt.findUser))
	rt.router.PUT("/settings/profile", rt.authMiddleware(rt.changeProfile))

	return rt.router
}
