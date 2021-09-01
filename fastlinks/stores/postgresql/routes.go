package postgresql

import "github.com/jmoiron/sqlx"

type RoutesStore struct {
	db *sqlx.DB
}

func NewRoutesStore(db *sqlx.DB) *RoutesStore {
	return &RoutesStore{
		db: db,
	}
}

func (rs *RoutesStore) Name() string {
	return "Routes Store"
}

func (rs *RoutesStore) Health() error {
	return rs.db.Ping()
}
