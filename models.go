package main

type User struct {
	Login    string
	Password string
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Text     string `json:"text,omitempty"`
	URL      string `json:"url,omitempty"`
	Category string `json:"category"`
}

type CreateCommentRequest struct {
	Comment string `json:"comment"`
}

type Comment struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

type Post struct {
	ID       int            `json:"id"`
	Title    string         `json:"title"`
	Text     string         `json:"text,omitempty"`
	URL      string         `json:"url,omitempty"`
	Author   string         `json:"author"`
	Category string         `json:"category"`
	Score    int            `json:"score"`
	Comments []Comment      `json:"comments"`
	Votes    map[string]int `json:"-"`
}
