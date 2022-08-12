package bbolt

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/TerrenceHo/monorepo/fastlinks/models"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
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

func TestPut(t *testing.T) {
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
		{
			route: &models.Route{
				Key:         "randKey",
				RedirectURL: "https://random.com",
			},
		},
	}

	dir := t.TempDir()

	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	for _, testcase := range testcases {
		err := routesStore.Put(testcase.route)
		if err != nil {
			assert.Fail("Add error: %s", err.Error())
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
		assert.Equalf(testcase.route, &route, "Added route not equivalent: %v != %v", testcase.route, route)
	}
}

func TestDelete(t *testing.T) {
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
		{
			route: &models.Route{
				Key:         "randKey",
				RedirectURL: "https://random.com",
			},
		},
	}

	dir := t.TempDir()
	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	for _, testcase := range testcases {
		err := routesStore.Put(testcase.route)
		if err != nil {
			assert.Fail("Add error: %s", err.Error())
		}

		// delete the route by key
		err = routesStore.Delete(testcase.route.Key)
		if err != nil {
			assert.Fail("Delete error: %s", err.Error())
		}

		// retrieve the route by key
		if err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(testBucket))
			v := b.Get([]byte(testcase.route.Key))
			if v != nil {
				return stackerrors.Errorf("deleted key should return nil data, returned %v", v)
			}
			return nil
		}); err != nil {
			assert.Nilf(err, "Error retrieving the stored route: %s", err.Error())
		}
	}
}

func TestGet(t *testing.T) {
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
		{
			route: &models.Route{
				Key:         "randKey",
				RedirectURL: "https://random.com",
			},
		},
	}

	dir := t.TempDir()
	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	for _, testcase := range testcases {
		err := routesStore.Put(testcase.route)
		if err != nil {
			assert.Fail("Add error: %s", err.Error())
		}
	}

	for _, testcase := range testcases {
		route, err := routesStore.GetByKey(testcase.route.Key)
		assert.Nil(err, "GetByKey error should be nil")
		assert.Equal(testcase.route, route, "route retrieved from Get not equal to expected")
	}
}

func TestGetAll(t *testing.T) {
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
		{
			route: &models.Route{
				Key:         "randKey",
				RedirectURL: "https://random.com",
			},
		},
	}

	dir := t.TempDir()
	db := getTestDB(t, dir)

	routesStore := NewRoutesStore(db, testBucket)
	routesStore.Migrate()

	for _, testcase := range testcases {
		err := routesStore.Put(testcase.route)
		if err != nil {
			assert.Fail("Add error: %s", err.Error())
		}
	}

	routes, err := routesStore.GetAll()
	if err != nil {
		assert.Failf("GetAll failed", "GetAll error: %s", err.Error())
	}
	assert.Equal(len(testcases), len(routes), "length of returned data not equal")
	for _, testcase := range testcases {
		found := false
		for _, route := range routes {
			if route.Key == testcase.route.Key {
				found = true
				break
			}
		}
		if !found {
			assert.Failf("GetAll validation failed", "could not find route with Key %s", testcase.route.Key)
		}
	}
}
