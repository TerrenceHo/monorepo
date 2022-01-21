package models

import (
	"testing"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/stretchr/testify/assert"
)

func TestRouteEncodeDecode(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		route *Route
	}

	testcases := []testcase{
		{
			route: &Route{
				Key: "key",
			},
		},
		{
			route: &Route{
				Key:         "key",
				RedirectURL: "https://github.com/TerrenceHo/monorepo",
			},
		},
		{
			route: &Route{
				Key:         "key",
				RedirectURL: "https://github.com/TerrenceHo/monorepo",
				ExtendedURL: "https://github.com/TerrenceHo/tree/master/{}",
			},
		},
	}

	for _, testcase := range testcases {
		bytes, err := testcase.route.Encode()
		if err != nil {
			t.Logf("encoding returned error: %s", err.Error())
		}
		var r Route
		err = r.Decode(bytes)
		if err != nil {
			t.Logf("decoding returned error: %s", err.Error())
		}
		assert.EqualValues(*testcase.route, r, "testcase route != decoded route")
	}
}

func validationError(msg string) error {
	return stackerrors.Wrap(
		stackerrors.New(msg),
		"validation failed",
	)
}

func TestValidation(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		route *Route
		err   error
	}

	testcases := []testcase{
		{
			route: &Route{
				Key:         "key1",
				RedirectURL: "https://github.com/TerrenceHo/monorepo",
			},
			err: nil,
		},
		{
			route: &Route{
				Key:         "key2",
				RedirectURL: "https://github.com/TerrenceHo/monorepo",
				ExtendedURL: "https://github.com/TerrenceHo/tree/master/{}",
			},
			err: nil,
		},
		{
			route: &Route{},
			err:   validationError("Route must have a Key, was empty string"),
		},
		{
			route: &Route{
				RedirectURL: "redirecturl",
			},
			err: validationError("Route must have a Key, was empty string"),
		},
		{
			route: &Route{
				Key: "key3",
			},
			err: validationError("Route must have a RedirectURL, was empty string"),
		},
		{
			route: &Route{
				Key:         "key4",
				RedirectURL: "https://github.com/TerrenceHo/monorepo",
				ExtendedURL: "https://github.com/TerrenceHo/tree/master/",
			},
			err: validationError("If route has ExtendedURL, it must contain '{}'"),
		},
	}

	for _, testcase := range testcases {
		err := testcase.route.Validate()
		if testcase.err == nil {
			if err != nil {
				assert.Nilf(err, "Validate error should be nil, error: %s", err.Error())
			}
		} else {
			assert.Equalf(testcase.err.Error(), err.Error(), "Validate error failed for key %s", testcase.route.Key)
		}
	}
}
