package services

import (
	"testing"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/stretchr/testify/assert"
)

type fakeStore struct{}

func (fs *fakeStore) Name() string  { return "fakeStore" }
func (fs *fakeStore) Health() error { return nil }

type fakeErrorStore struct{}

func (fes *fakeErrorStore) Name() string { return "fakeErrorStore" }
func (fes *fakeErrorStore) Health() error {
	return stackerrors.New("fake store error")
}

func TestCheckHealth(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		storesList []HealthStore
		want       []HealthCheckError
	}

	testcases := []testcase{
		{
			storesList: []HealthStore{
				&fakeStore{},
			},
			want: []HealthCheckError{},
		},
		{
			storesList: []HealthStore{
				&fakeErrorStore{},
			},
			want: []HealthCheckError{
				{
					Name:  "fakeErrorStore",
					Error: stackerrors.New("fake store error"),
				},
			},
		},
		{
			storesList: []HealthStore{
				&fakeStore{},
				&fakeErrorStore{},
			},
			want: []HealthCheckError{
				{
					Name:  "fakeErrorStore",
					Error: stackerrors.New("fake store error"),
				},
			},
		},
	}

	for _, testcase := range testcases {
		hs := NewHealthService(testcase.storesList...)
		errorsList := hs.Check()
		assert.Lenf(
			errorsList, len(testcase.want),
			"Lengths of errors not equal: %v != %v", errorsList, testcase.want,
		)
		for i, err := range errorsList {
			assert.EqualValuesf(
				testcase.want[i].Name, err.Name,
				"names not equal: %s != %s", testcase.want[i].Name, err.Name,
			)
			assert.EqualValuesf(
				testcase.want[i].Error.Error(), err.Error.Error(),
				"errors not equal: %s != %s",
				testcase.want[i].Error.Error(), err.Error.Error(),
			)
		}
	}
}
