package main

// User entity
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// Todo entity
type Todo struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	DateCreated string `json:"dateCreated"`
}

// ErrorResponse is used to return a JSON response on server error
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
