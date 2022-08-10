package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TerrenceHo/monorepo/fastlinks"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	logo = `
 ______   ______     ______     ______   __         __     __   __     __  __     ______    
/\  ___\ /\  __ \   /\  ___\   /\__  _\ /\ \       /\ \   /\ "-.\ \   /\ \/ /    /\  ___\   
\ \  __\ \ \  __ \  \ \___  \  \/_/\ \/ \ \ \____  \ \ \  \ \ \-.  \  \ \  _"-.  \ \___  \  
 \ \_\    \ \_\ \_\  \/\_____\    \ \_\  \ \_____\  \ \_\  \ \_\\"\_\  \ \_\ \_\  \/\_____\ 
  \/_/     \/_/\/_/   \/_____/     \/_/   \/_____/   \/_/   \/_/ \/_/   \/_/\/_/   \/_____/
`
)

func rootCmd(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	mainCmd := &cobra.Command{
		Use:   "fastlinks",
		Short: "fastlinks: simple URL redirector",
		Run:   run,
		Args:  cobra.ArbitraryArgs,
	}

	flags := mainCmd.Flags()

	flags.String("env", "dev", "environment: test, dev, or prod")
	flags.Bool("hidebanner", false, "hide fastlinks banner")
	flags.String("host", "localhost", "hostserver on this hostname")
	flags.StringP("port", "p", "12345", "host server on localhost:<port>")

	flags.String("db.user", "fastlinks", "database user")
	flags.String("db.password", "password", "database password")
	flags.String("db.dbname", "fastlinks", "database name")
	flags.String("db.port", "5432", "database port")
	flags.String("db.host", "localhost", "database host")
	flags.String("db.sslmode", "disable", "database sslmode")

	return mainCmd
}

func main() {
	mainCmd := rootCmd(serve)
	mainCmd.Execute()
}

func serve(c *cobra.Command, args []string) {
	conf, errs := loadConfig(c)
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err.Error())
		}
		os.Exit(1)
	}
	if !conf.Hidebanner {
		fmt.Println(logo)
	}

	fastlinks.Start(conf)
}

// loadConfig takes in a cobra Command and unmarshalls it into the custom config
// struct. It also binds pflags to viper, sets environment variables, and reads
// in a fastlinks-config.{json,yaml,toml}.
func loadConfig(cmd *cobra.Command) (fastlinks.Config, []error) {
	var errs []error
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		errs = append(errs, stackerrors.Wrap(err, "failed to bind flags"))
	}

	var conf fastlinks.Config
	if err := viper.Unmarshal(&conf); err != nil {
		errs = append(errs, stackerrors.Wrap(err, "failed to unmarshal config"))
	}
	return conf, errs
}
