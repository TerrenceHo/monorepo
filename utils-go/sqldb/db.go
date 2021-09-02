package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Direction string

const (
	Up   Direction = "up"
	Down Direction = "down"
)

// NewConnection instantiates a new connection, returning a sqlx.DB instance.
// Currently, it accepts postgres as an engine.
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

// DB is a wrapper around *.sqlx.DB, adding some migration tooling to the
// database instance.
type DB struct {
	*sqlx.DB
}

// NewSQLDB takes in a *sqlx.DB instance and returns a *DB instance.
func NewSQLDB(db *sqlx.DB) *DB {
	return &DB{
		db,
	}
}

func (db *DB) SQLxDB() *sqlx.DB {
	return db.DB
}

func (db *DB) SQLDB() *sql.DB {
	return db.DB.DB
}

func (db *DB) Migrate(direction Direction, sourceURL string) error {
	driver, err := postgres.WithInstance(db.SQLDB(), &postgres.Config{})
	if err != nil {
		return stackerrors.Wrap(err, "failed to create database migration driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return stackerrors.Wrap(err, "failed to create database migration instance")
	}
	switch direction {
	case Up:
		if err = m.Up(); err != nil {
			return stackerrors.Wrapf(err, "failed to migrate upwards for migration %s", sourceURL)
		}
	case Down:
		if err = m.Down(); err != nil {
			return stackerrors.Wrapf(err, "failed to migrate downwards for migration %s", sourceURL)
		}

	default:
		return stackerrors.New("migration direction must be either up or down")
	}
	return nil
}
