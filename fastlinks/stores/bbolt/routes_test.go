package bbolt

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/TerrenceHo/monorepo/fastlinks/models"
	"github.com/stretchr/testify/assert"

	bolt "go.etcd.io/bbolt"
)

const testBucket = "test-bucket"

func getTestDB(t *testing.T, dir string) *bolt.DB {
	file, err := ioutil.TempFile(dir, "test.db")
	db, err := NewConnection(file.Name(), 0600)
	if err != nil {
		t.Fatalf("failed to open bbolt db: %s", err.Error())
	}
	return db
}

func TestMigrate(t *testing.T) {
	assert := assert.New(t)

	dir := t.TempDir()

	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(testBucket))

		// if b is nil, then bucket does not exist
		assert.NotNilf(b, "bucket %s should not be nil", testBucket)
		return nil
	})

	if err != nil {
		assert.Nilf(err, "error during migration: %s", err.Error())
	}
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		route *models.Route
	}

	testcases := []testcase{
		{
			route: &models.Route{
				Key:         "testkey",
				RedirectURL: "https://google.com",
				ExtendedURL: "https://google.com/search/{}",
			},
		},
	}

	dir := t.TempDir()

	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	for _, testcase := range testcases {
		err := routesStore.Create(testcase.route)
		if err != nil {
			assert.Fail("Create error: %s", err.Error())
		}

		// retrieve the route by key
		var route models.Route
		if err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(testBucket))
			routeBytes := bytes.NewBuffer(b.Get([]byte(testcase.route.Key)))
			err := (&route).Decode(routeBytes)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			assert.Nilf(err, "Error retrieving the stored route: %s", err.Error())
		}
		assert.Equalf(testcase.route, &route, "Created route not equivalent: %v != %v", testcase.route, route)
	}
}
