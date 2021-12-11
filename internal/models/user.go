package models

type User struct {
	Name  string `json:"name"`
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
