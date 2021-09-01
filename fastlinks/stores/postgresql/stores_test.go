package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectPostgresql(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		user     string
		password string
		dbname   string
		port     string
		host     string
		sslmode  string
		want     string
	}

	testcases := []testcase{
		{
			user:     "dnd",
			password: "barbarian",
			dbname:   "5e",
			port:     "5357",
			host:     "localhost",
			sslmode:  "disable",
			want:     "user=dnd password=barbarian dbname=5e port=5357 host=localhost sslmode=disable",
		},
	}

	for _, testcase := range testcases {
		dbConnection := connectPostgresql(
			testcase.user,
			testcase.password,
			testcase.dbname,
			testcase.port,
			testcase.host,
			testcase.sslmode,
		)
		assert.EqualValuesf(
			testcase.want, dbConnection, "%s != %s", testcase.want, dbConnection,
		)
	}
}
