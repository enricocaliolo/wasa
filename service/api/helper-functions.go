package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getToken(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	token := parts[1]
	num, err := strconv.Atoi(token)
	if err != nil {
		log.Fatal(err)
	}

	return num
}
