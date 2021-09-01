package sqldb

import (
	"fmt"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection(
	engine, user, password, dbname, port, host, sslmode string,
) (*sqlx.DB, error) {
	var dbConnection string
	switch engine {
	case "postgres":
		dbConnection = connectPostgresql(user, password, dbname, port, host, sslmode)
	default:
		dbConnection = ""
	}

	db, err := sqlx.Connect(engine, dbConnection)
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

type DB struct {
	*sqlx.DB
}

func NewSQLDB(db *sqlx.DB) *DB {
	return &DB{
		db,
	}
}
