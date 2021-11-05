package repository

import (
	"github.com/nentenpizza/citadels/internal/domain"
	"github.com/nentenpizza/citadels/internal/repository/postgres"
)

type Users interface {
	Create(string) (int64, error)
	ByID(int64) (domain.User, error)
	ByName(string) (domain.User, error)
	ExistsByName(string) (bool, error)
}

type Repositories struct {
	Users Users
}

func New(db *postgres.DB) *Repositories{
	return &Repositories{
		Users: postgres.NewUsersRepository(db),
	}
}