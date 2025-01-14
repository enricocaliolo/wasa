package api

import (
	"database/sql"
	"encoding/json"
	"io"
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

	user, err := rt.db.GetUser(username)
	if err == sql.ErrNoRows {
		user, err = rt.db.CreateUser(username)
		if err != nil {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)

}

func (rt *APIRouter) GetAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users, err := rt.db.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Error getting users")
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(users)

}

func (rt *APIRouter) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query()

	username := query.Get("username")

	user, err := rt.db.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("User not found")
		return
	}

	w.Header().Set("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)

}

func (rt *APIRouter) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user models.User
	user.ID = userId

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isUserUpdated, err := rt.db.UpdateUsername(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Error getting users")
		return
	}

	if !isUserUpdated {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode("Username already taken!")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *APIRouter) changePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isUserUpdated, err := rt.db.UpdatePhoto(userId, imageData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Error getting users")
		return
	}

	if !isUserUpdated {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Username already taken!")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Icon updated succesfully")
}
