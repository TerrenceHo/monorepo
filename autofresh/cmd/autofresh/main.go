package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/TerrenceHo/monorepo/autofresh"
	pth "github.com/TerrenceHo/monorepo/utils-go/path"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rootCmd(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	mainCmd := &cobra.Command{
		Use:   "autofresh",
		Short: "Autofresh: live reloading server",
		Run:   run,
		Args:  cobra.ArbitraryArgs,
	}

	flags := mainCmd.Flags()

	flags.StringSliceP("extensions", "e", []string{}, "comma separated string of file extensions")
	flags.Bool("hidebanner", false, "hide banner")
	flags.StringSliceP("watch", "w", []string{}, "watch specific directory")

	return mainCmd
}

func main() {
	mainCmd := rootCmd(serve)
	mainCmd.Execute()
}

func serve(cmd *cobra.Command, args []string) {
	c, errs := loadConfig(cmd)
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err.Error())
		}
		os.Exit(1)
	}
	autofresh.Start(purifyConfig(c, args))
}

// purifyConfig takes in the config object and the leftover arguments from the
// cobra command object and does the following:
// (1) Sets the Cmd and Args variables in the config.
// (2) If Watch is not specified, adds the current directory to watch.
// (3) Resolve all watch directories to absolute paths, after expanding user
// home directories.
func purifyConfig(c autofresh.Config, args []string) autofresh.Config {
	c.Cmd = args[0]
	c.Args = args[1:]
	if len(c.Watch) == 0 {
		path, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		c.Watch = []string{path}
	}

	for i := range c.Watch {
		abspath, err := filepath.Abs(pth.ExpandUser(c.Watch[i]))
		if err != nil {
			log.Fatalf(
				"Failed to get absolute path for path (%s): %v", abspath, err,
			)
		}
		c.Watch[i] = abspath
	}
	return c
}

// LoadConfig takes in a cobra Command and unmarshalls it into the custom config
// struct. It also binds pflags to viper, sets environment variables, and reads
// in a autofresh-config.{json,yaml,toml}.
func loadConfig(cmd *cobra.Command) (autofresh.Config, []error) {
	var errs []error
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		errs = append(errs, stackerrors.Wrap(err, "failed to bind flags"))
	}

	var conf autofresh.Config
	if err := viper.Unmarshal(&conf); err != nil {
		errs = append(errs, stackerrors.Wrap(err, "failed to unmarshal config"))
	}
	return conf, errs
}
