package postgres

import "github.com/nentenpizza/citadels/internal/domain"

type UsersRepository struct {
	db *DB
}

func NewUsersRepository(db *DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u UsersRepository) Create(s string) (int64, error) {
	panic("implement me")
}

func (u UsersRepository) ByID(i int64) (domain.User, error) {
	panic("implement me")
}

func (u UsersRepository) ByName(s string) (domain.User, error) {
	panic("implement me")
}

func (u UsersRepository) ExistsByName(s string) (bool, error) {
	panic("implement me")
}
