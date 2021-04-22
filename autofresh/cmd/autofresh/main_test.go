package main

import (
	"testing"

	"github.com/TerrenceHo/monorepo/autofresh"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func generateConfig(t *testing.T, args []string) autofresh.Config {
	var c autofresh.Config
	testServe := func(cmd *cobra.Command, args []string) {
		var errs []error
		c, errs = loadConfig(cmd)

		if len(errs) != 0 {
			for _, err := range errs {
				t.Log(err.Error())
			}
			t.Fatal("load config failed")
		}
		c.Cmd = args
	}
	cmd := rootCmd(testServe)
	cmd.SetArgs(args)
	cmd.Execute()
	return c
}

func TestRootCmdConfig(t *testing.T) {
	assert := assert.New(t)
	type testcase struct {
		args []string
		want autofresh.Config
	}
	testcases := []testcase{
		{
			args: []string{
				"--watch", "/User",
				"-w", "~/go",
				"--hidebanner",
				"-w", "github.com/monorepo",
				"-e", "go,js,ts",
				"--", "bash", "--help",
			},
			want: autofresh.Config{
				Cmd:        []string{"bash", "--help"},
				Extensions: []string{"go", "js", "ts"},
				HideBanner: true,
				Watch:      []string{"/User", "~/go", "github.com/monorepo"},
			},
		},
	}
	for _, testcase := range testcases {
		c := generateConfig(t, testcase.args)
		assert.EqualValuesf(
			testcase.want, c,
			"configs are not equal: %v != %v", testcase.want, c,
		)
	}
}
