package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterUserHandlers(s *Storage) {
	http.HandleFunc("/api/user/", userHandler(s))
}

func userHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		login := strings.TrimPrefix(r.URL.Path, "/api/user/")
		json.NewEncoder(w).Encode(s.GetPostsByUser(login))
	}
}
