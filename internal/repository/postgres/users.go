package postgres

import "github.com/nentenpizza/citadels/internal/models"

type UsersRepository struct {
	db *DB
}

func NewUsersRepository(db *DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u UsersRepository) Create(name string, email string, passwordHash string) (int64, error) {
	const q = `insert into users (name, email, password_hash) values ($1, $2, $3) returning id`
	rows, err := u.db.Query(q, name, email, passwordHash)
	if err != nil {
		return 0, err
	}
	var id int64
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return id, err
		}
	}
	return id, nil
}

func (u UsersRepository) ByID(id int64) (user models.User, _ error) {
	const q = `select * from users where id = $1`
	return user, u.db.Get(&user, q, id)
}

func (u UsersRepository) ByName(name string) (user models.User, _ error) {
	const q = `select * from users where name = $1`
	return user, u.db.Get(&user, q, name)
}

func (u UsersRepository) ExistsByName(name string) (has bool, _ error) {
	const q = `select exists(*) from users where name = $1`
	return has, u.db.Get(&has, q, name)
}
