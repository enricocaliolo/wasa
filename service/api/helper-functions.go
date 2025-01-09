package api

import (
	"net/http"
	"strconv"
	"strings"
)

func getToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	token := parts[1]
	num, err := strconv.Atoi(token)
	if err != nil {
		return -1, nil
	}

	return num, nil
}
