package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) TestUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users := rt.db.GetAllUsers()
	print(users)

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(users)

}
