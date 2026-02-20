package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterAuthHandlers(storage *Storage) {
	http.HandleFunc("/api/register", registerHandler(storage))
	http.HandleFunc("/api/login", loginHandler(storage))
}

func registerHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Login == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := storage.Register(req.Login, req.Password); err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
	}
}

func loginHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := s.Login(req.Login, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func getUserFromRequest(r *http.Request, s *Storage) (string, bool) {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return "", false
	}
	return s.GetUserByToken(strings.TrimPrefix(token, "Bearer "))
}
