package bbolt

import (
	"os"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	bolt "go.etcd.io/bbolt"
)

type Store interface{}

func NewConnection(filepath string, mode os.FileMode) (*bolt.DB, error) {
	db, err := bolt.Open(filepath, mode, nil)
	if err != nil {
		return nil, stackerrors.Wrap(err, "bbolt DB failed to open")
	}
	return db, nil
}
