package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
)

func (rt *APIRouter) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// type req struct {
	// 	username string
	// }

	var req struct {
		Username string `json:"username"`
	 }
	 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	 }
	 username := req.Username

	_, err := rt.db.GetUser(username)
	if err == sql.ErrNoRows {
		_, err = rt.db.CreateUser(username)
		if err != nil {
			return
		}
	}
	resp := struct {
        Username string `json:"username"`
    }{
        Username: req.Username,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)

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

func (rt *APIRouter) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	user.ID = getToken(r)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isUserUpdated := rt.db.UpdateUsername(user)

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

func (rt *APIRouter) changePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	user.ID = getToken(r)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isUserUpdated := rt.db.UpdatePhoto(user)

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
