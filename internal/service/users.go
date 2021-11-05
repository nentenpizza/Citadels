package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/nentenpizza/citadels/internal/repository"
)

var (
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

type UserRegisterForm struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"min=3,max=20"`
	Password string `json:"password" validate:"min=6"`
}

type Users struct {
	Repos *repository.Repositories
}

func (s Users) Register(form UserRegisterForm) error {
	exists, err := s.Repos.Users.ExistsByName(form.Username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyRegistered
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.Repos.Users.Create(form.Username, form.Email, string(hash))
	return err
}
