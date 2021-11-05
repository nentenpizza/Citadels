package v1

import "github.com/nentenpizza/citadels/internal/service"

type Handler struct {
	Users service.Users
}

func NewHandler(users service.Users) *Handler {
	return &Handler{Users: users}
}
