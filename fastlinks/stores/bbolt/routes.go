package bbolt

import (
	"bytes"

	"github.com/TerrenceHo/monorepo/fastlinks/models"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"

	bolt "go.etcd.io/bbolt"
)

type RoutesStore struct {
	db     *bolt.DB
	bucket string
}

func NewRoutesStore(db *bolt.DB, bucketName string) *RoutesStore {
	return &RoutesStore{
		db:     db,
		bucket: bucketName,
	}
}

func (rs *RoutesStore) Migrate() error {
	if err := rs.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(rs.bucketName()))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return stackerrors.Wrap(err, "failed to create bbolt bucket routes")
	}

	return nil
}

func (rs *RoutesStore) Name() string {
	return "BBolt Routes Store"
}

func (rs *RoutesStore) Health() error {
	return nil
}

func (rs *RoutesStore) GetByKey(key string) (*models.Route, error) {
	var route models.Route
	if err := rs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(rs.bucketName())
		routeBytes := bytes.NewBuffer(b.Get([]byte(key)))
		err := (&route).Decode(routeBytes)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, stackerrors.Wrapf(err, "failed to get route with key %s", key)
	}
	return &route, nil
}

func (rs *RoutesStore) GetAll() ([]*models.Route, error) {
	var routes []*models.Route

	if err := rs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(rs.bucketName())
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var route models.Route
			routeBytes := bytes.NewBuffer(v)
			err := (&route).Decode(routeBytes)
			if err != nil {
				return err
			}
			routes = append(routes, &route)
		}
		return nil
	}); err != nil {
		return nil, stackerrors.Wrap(err, "failed to query all routes")
	}

	return routes, nil
}

func (rs *RoutesStore) Put(route *models.Route) error {
	if err := rs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(rs.bucketName())
		routeBytes, err := route.Encode()
		if err != nil {
			return err
		}
		return b.Put([]byte(route.Key), routeBytes.Bytes())
	}); err != nil {
		return stackerrors.Wrapf(err, "failed to put route with key %s", route.Key)
	}
	return nil
}

func (rs *RoutesStore) Delete(key string) error {
	if err := rs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(rs.bucketName())
		return b.Delete([]byte(key))
	}); err != nil {
		return stackerrors.Wrapf(err, "failed to delete route with key %s", key)
	}
	return nil
}

func (rs *RoutesStore) bucketName() []byte {
	return []byte(rs.bucket)
}
