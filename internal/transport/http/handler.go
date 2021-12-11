package httpapi

import "github.com/nentenpizza/citadels/internal/service"

type AuthHandler struct {
	Users service.Users
}

func NewAuthHandler(users service.Users) *AuthHandler {
	return &AuthHandler{Users: users}
}
