package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query()

	username := query.Get("username")

	user := rt.db.GetUser(username)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if user.ID == -1 {
		json.NewEncoder(w).Encode("User Not Found")
		return
	}

	json.NewEncoder(w).Encode(user)

}
