package postgres

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Open(url string) (*DB, error) {
	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return &DB{
		DB: db,
	}, nil
}
