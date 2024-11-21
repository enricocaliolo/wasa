package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *APIRouter) authMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// validating
		user := rt.db.GetUser(token)
		if user.ID == -1 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	}
}
