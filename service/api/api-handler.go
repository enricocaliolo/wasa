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
	rt.router.PUT("/settings/profile/username", rt.authMiddleware(rt.changeUsername))
	rt.router.PUT("/settings/profile/photo", rt.authMiddleware(rt.changePhoto))

	// conversation operations
	rt.router.GET("/conversations", rt.authMiddleware(rt.conversations))
	rt.router.GET("/conversations/:id", rt.authMiddleware(rt.getConversation))
	rt.router.POST("/conversations/:id", rt.authMiddleware(rt.sendMessage))
	rt.router.DELETE("/conversations/:id", rt.authMiddleware(rt.deleteConversation))
	rt.router.POST("/conversations/:id/reply", rt.authMiddleware(rt.sendMessage))
	rt.router.POST("/conversations/:id/forward", rt.authMiddleware(rt.forwardMessage))

	// rt.router.POST("/conversations")
	// rt.router.PUT("/conversations/:id/users")
	// rt.router.PUT("/conversations/:id/name")
	// rt.router.PUT("/conversations/:id/photo")
	// rt.router.PUT("/conversations/{conversation_id}/messages/{message_id}")
	// rt.router.DELETE("/conversations/{conversation_id}/messages/{message_id}/reactions/{reaction_id}")

	return rt.router
}
