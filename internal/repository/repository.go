package repository

import (
	"github.com/nentenpizza/citadels/internal/models"
	"github.com/nentenpizza/citadels/internal/repository/postgres"
)

type Users interface {
	Create(name string, email string, password string) (int64, error)
	ByID(int64) (models.User, error)
	ByName(string) (models.User, error)
	ExistsByName(string) (bool, error)
}

type Repositories struct {
	Users Users
}

func New(db *postgres.DB) *Repositories {
	return &Repositories{
		Users: postgres.NewUsersRepository(db),
	}
}
