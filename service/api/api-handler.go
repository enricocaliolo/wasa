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

	// conversation operations
	rt.router.GET("/conversations", rt.authMiddleware(rt.conversations))
	rt.router.GET("/conversations/:id", rt.authMiddleware(rt.getConversation))
	rt.router.POST("/conversations/:id", rt.authMiddleware(rt.sendMessage))
	rt.router.DELETE("/conversations/:id", rt.authMiddleware(rt.deleteConversation))
	rt.router.POST("/conversations/:id/reply", rt.authMiddleware(rt.sendMessage))
	rt.router.POST("/conversations/:id/forward", rt.authMiddleware(rt.forwardMessage))

	return rt.router
}
