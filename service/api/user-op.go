package api

import (
	"encoding/json"
	"net/http"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
)

func (rt *APIRouter) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	w.Header().Set("Authorization", "Bearer "+string(rune(user.ID)))
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func (rt *APIRouter) TestUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users := rt.db.GetAllUsers()
	print(users)

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(users)

}

func (rt *APIRouter) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (rt *APIRouter) changeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	user.ID = getToken(r)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userUpdated := rt.db.UpdateProfile(user)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userUpdated)

}
