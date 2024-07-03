package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}
	err2 := db.Ping()
	if err2 != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return &Store{
		ThreadStore:  &ThreadStore{DB: db},
		PostStore:    &PostStore{DB: db},
		CommentStore: &CommentStore{DB: db},
	}, nil
}

type Store struct {
	*ThreadStore
	*PostStore
	*CommentStore
}
