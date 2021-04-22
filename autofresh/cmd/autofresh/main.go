package main

import (
	"fmt"
	"os"

	"github.com/TerrenceHo/monorepo/autofresh"
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
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	c.Cmd = args
	autofresh.Start(c)
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
