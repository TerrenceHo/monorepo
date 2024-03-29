package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TerrenceHo/monorepo/fastlinks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func generateConfig(t *testing.T, args []string) fastlinks.Config {
	var c fastlinks.Config
	testServe := func(cmd *cobra.Command, args []string) {
		var errs []error
		c, errs = loadConfig(cmd)

		if len(errs) != 0 {
			for _, err := range errs {
				t.Log(err.Error())
			}
			t.Fatal("load config failed")
		}
	}
	cmd := rootCmd(testServe)
	cmd.SetArgs(args)
	cmd.Execute()
	return c
}

func TestConfig(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	localFile := filepath.Join(home, ".fastlinks/fastlinks.db")

	assert := assert.New(t)
	type testcase struct {
		args []string
		want fastlinks.Config
	}
	testcases := []testcase{
		{
			args: []string{},
			want: fastlinks.Config{
				Env:        "dev",
				Hidebanner: false,
				Host:       "localhost",
				Port:       "12345",
				Storage:    "local",
				DB: fastlinks.DBConfig{
					User:     "fastlinks",
					Password: "password",
					DBName:   "fastlinks",
					Port:     "5432",
					Host:     "localhost",
					SSLMode:  "disable",
				},
				Local: fastlinks.LocalConfig{
					File: localFile,
				},
			},
		},
		{
			args: []string{
				"--env=prod",
				"--hidebanner",
				"--host=google.com",
				"--port=5555",
				"--db.user=user",
				"--db.password=newpassword",
				"--db.dbname=newdb",
				"--db.port=6666",
				"--db.host=newhost.com",
				"--db.sslmode=verify-full",
			},
			want: fastlinks.Config{
				Env:        "prod",
				Hidebanner: true,
				Host:       "google.com",
				Port:       "5555",
				Storage:    "local",
				DB: fastlinks.DBConfig{
					User:     "user",
					Password: "newpassword",
					DBName:   "newdb",
					Port:     "6666",
					Host:     "newhost.com",
					SSLMode:  "verify-full",
				},
				Local: fastlinks.LocalConfig{
					File: localFile,
				},
			},
		},
	}

	for _, testcase := range testcases {
		c := generateConfig(t, testcase.args)
		assert.EqualValuesf(
			testcase.want, c, "configs are not equal: %v != %v", testcase.want, c,
		)
	}
}
