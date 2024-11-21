package api

import (
	"encoding/json"
	"log"
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

	id, err := rt.db.GetUser(username)
	if err != nil {
		log.Fatal(err)
		return
	}

	if id == -1 {
		id, err = rt.db.CreateUser(username)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	w.Header().Set("Authorization", "Bearer "+string(rune(id)))
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(id)

}

func (rt *APIRouter) TestUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users := rt.db.GetAllUsers()

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(users)

}

func (rt *APIRouter) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query()

	username := query.Get("username")

	id, err := rt.db.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not found")
		return
	}

	w.Header().Set("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)

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

	isUserUpdated := rt.db.UpdateProfile(user)

	if !isUserUpdated {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Username already taken!")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
