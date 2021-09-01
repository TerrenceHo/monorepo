package postgresql

import (
	"github.com/TerrenceHo/monorepo/utils-go/sqldb"
)

type RoutesStore struct {
	db *sqldb.DB
}

func NewRoutesStore(db *sqldb.DB) *RoutesStore {
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
