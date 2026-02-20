package main

import (
	"errors"
	"sync"
)

type Storage struct {
	mu         sync.RWMutex
	Users      map[string]User
	Tokens     map[string]string
	Posts      []Post
	NextPostID int
	NextComID  int
}

func NewStorage() *Storage {
	return &Storage{
		Users:      make(map[string]User),
		Tokens:     make(map[string]string),
		Posts:      []Post{},
		NextPostID: 1,
		NextComID:  1,
	}
}

func (s *Storage) Register(login, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Users[login]; exists {
		return errors.New("user exists")
	}

	s.Users[login] = User{Login: login, Password: password}
	return nil
}

func (s *Storage) Login(login, password string) (string, error) {
	s.mu.RLock()
	user, exists := s.Users[login]
	s.mu.RUnlock()

	if !exists || user.Password != password {
		return "", errors.New("invalid credentials")
	}

	token := login + "_token"

	s.mu.Lock()
	s.Tokens[token] = login
	s.mu.Unlock()

	return token, nil
}

func (s *Storage) GetUserByToken(token string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	login, ok := s.Tokens[token]
	return login, ok
}

func (s *Storage) CreatePost(login string, req CreatePostRequest) (Post, error) {
	if req.Title == "" || req.Category == "" {
		return Post{}, errors.New("invalid")
	}

	if req.Text == "" && req.URL == "" {
		return Post{}, errors.New("empty content")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	post := Post{
		ID:       s.NextPostID,
		Title:    req.Title,
		Text:     req.Text,
		URL:      req.URL,
		Author:   login,
		Category: req.Category,
		Votes:    make(map[string]int),
	}

	s.NextPostID++
	s.Posts = append(s.Posts, post)

	return post, nil
}

func (s *Storage) GetPosts(category string) []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if category == "" {
		return s.Posts
	}

	var result []Post
	for _, p := range s.Posts {
		if p.Category == category {
			result = append(result, p)
		}
	}
	return result
}

func (s *Storage) GetPost(id int) (*Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Posts {
		if s.Posts[i].ID == id {
			return &s.Posts[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (s *Storage) AddComment(postID int, login, body string) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Posts {
		if s.Posts[i].ID == postID {
			comment := Comment{
				ID:     s.NextComID,
				Author: login,
				Body:   body,
			}
			s.NextComID++
			s.Posts[i].Comments = append(s.Posts[i].Comments, comment)
			return &s.Posts[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (s *Storage) DeleteComment(postID, commentID int, login string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Posts {
		if s.Posts[i].ID == postID {
			for j, c := range s.Posts[i].Comments {
				if c.ID == commentID {
					if c.Author != login {
						return errors.New("forbidden")
					}
					s.Posts[i].Comments = append(
						s.Posts[i].Comments[:j],
						s.Posts[i].Comments[j+1:]...,
					)
					return nil
				}
			}
		}
	}
	return errors.New("not found")
}

func (s *Storage) DeletePost(postID int, login string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.Posts {
		if p.ID == postID {
			if p.Author != login {
				return errors.New("forbidden")
			}
			s.Posts = append(s.Posts[:i], s.Posts[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (s *Storage) Vote(postID int, login string, value int) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Posts {
		if s.Posts[i].ID == postID {

			prev := s.Posts[i].Votes[login]
			s.Posts[i].Score -= prev

			if value == 0 {
				delete(s.Posts[i].Votes, login)
			} else {
				s.Posts[i].Votes[login] = value
				s.Posts[i].Score += value
			}

			return &s.Posts[i], nil
		}
	}
	return nil, errors.New("not found")
}
func (s *Storage) GetPostsByUser(login string) []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Post
	for _, p := range s.Posts {
		if p.Author == login {
			result = append(result, p)
		}
	}
	return result
}
