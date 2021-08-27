package main

import (
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
	assert := assert.New(t)
	type testcase struct {
		args []string
		want fastlinks.Config
	}
	testcases := []testcase{
		{
			args: []string{
				"--hidebanner",
			},
			want: fastlinks.Config{
				Hidebanner: true,
				Port:       12345,
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
