package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RegisterPostHandlers(storage *Storage) {
	http.HandleFunc("/api/posts/", postsHandler(storage))
	http.HandleFunc("/api/post/", postHandler(storage))
}

func postsHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		path = strings.Trim(path, "/")

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(s.GetPosts(path))

		case http.MethodPost:
			login, ok := getUserFromRequest(r, s)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var req CreatePostRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			post, err := s.CreatePost(login, req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(post)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/post/"), "/")

		if len(parts) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		login, _ := getUserFromRequest(r, s)

		if len(parts) == 2 {
			action := parts[1]

			switch action {

			case "upvote":
				if login == "" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				post, err := s.Vote(id, login, 1)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(post)
				return

			case "downvote":
				if login == "" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				post, err := s.Vote(id, login, -1)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(post)
				return

			case "unvote":
				if login == "" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				post, err := s.Vote(id, login, 0)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(post)
				return
			}
		}

		switch r.Method {

		case http.MethodGet:
			post, err := s.GetPost(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(post)

		case http.MethodPost:
			if login == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var req CreateCommentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			post, err := s.AddComment(id, login, req.Comment)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(post)

		case http.MethodDelete:
			if len(parts) == 2 {
				commentID, err := strconv.Atoi(parts[1])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if err := s.DeleteComment(id, commentID, login); err != nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				return
			}

			if err := s.DeletePost(id, login); err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
