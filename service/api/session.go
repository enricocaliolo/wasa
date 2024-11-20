package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) session(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user := rt.db.GetUser(username)

	if user.ID == -1 {
		user = rt.db.CreateUser(username)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(user)

}
