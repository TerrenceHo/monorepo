package models

import (
	"testing"

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
