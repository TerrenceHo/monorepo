package postgresql

import (
	"fmt"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store interface{}

func NewConnection(
	user, password, dbname, port, host, sslmode string,
) (*sqlx.DB, error) {
	dbConnection := connectPostgresql(user, password, dbname, port, host, sslmode)

	db, err := sqlx.Connect("postgres", dbConnection)
	if err != nil {
		return nil, stackerrors.Wrap(err, "postgresql connection failed to connect")
	}

	return db, nil
}

func connectPostgresql(user, password, dbname, port, host, sslmode string) string {
	dbConnection := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s",
		user, password, dbname, port, host, sslmode)
	return dbConnection
}
