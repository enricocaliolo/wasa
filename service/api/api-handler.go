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
	rt.router.PUT("/settings/profile/icon", rt.authMiddleware(rt.changePhoto))

	// conversation operations
	rt.router.GET("/conversations", rt.authMiddleware(rt.conversations))
	rt.router.GET("/conversations/:conversation_id", rt.authMiddleware(rt.getConversation))
	rt.router.POST("/conversations/:conversation_id", rt.authMiddleware(rt.sendMessage))
	rt.router.DELETE("/conversations/:conversation_id", rt.authMiddleware(rt.deleteConversation))
	rt.router.POST("/conversations/:conversation_id/reply", rt.authMiddleware(rt.sendMessage))
	rt.router.POST("/conversations/:conversation_id/forward", rt.authMiddleware(rt.sendMessage))
	rt.router.DELETE("/conversations/:conversation_id/messages/:message_id", rt.authMiddleware(rt.deleteMessage))

	rt.router.PUT("/conversations/:conversation_id/messages/:message_id", rt.authMiddleware(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversation_id/messages/:message_id/reactions/:reaction_id", rt.authMiddleware(rt.uncommentMessage))

	rt.router.POST("/conversations", rt.authMiddleware(rt.createConversation))
	rt.router.PUT("/conversations/:conversation_id/name", rt.authMiddleware(rt.updateGroupName))
	rt.router.PUT("/conversations/:conversation_id/photo", rt.authMiddleware(rt.UpdateGroupPhoto))
	rt.router.PUT("/conversations/:conversation_id/users", rt.authMiddleware(rt.addGroupMembers))

	rt.router.ServeFiles("/files/*filepath", http.Dir("files"))
	rt.router.GET("/ws", rt.authMiddleware(rt.HandleWebSocket))

	return rt.router
}
