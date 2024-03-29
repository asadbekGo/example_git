package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"app/config"
	"app/storage"
)

type Store struct {
	db   *pgxpool.Pool
	book storage.BookRepoI
}

func NewConnectPostgresql(cfg *config.Config) (storage.StorageI, error) {

	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))
	if err != nil {
		return nil, err
	}

	pgpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:   pgpool,
		book: NewBookRepo(pgpool),
	}, nil
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) Book() storage.BookRepoI {

	if s.book == nil {
		s.book = NewBookRepo(s.db)
	}

	return s.book
}

// GORM
// ROW
// SQLBUILDER
// SQLX
// PGXPOOL
