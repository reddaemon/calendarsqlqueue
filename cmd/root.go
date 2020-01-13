package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"

	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendarsqlqueue/app"
	"github.com/reddaemon/calendarsqlqueue/config"
	database "github.com/reddaemon/calendarsqlqueue/db"
	"github.com/reddaemon/calendarsqlqueue/queue"
	"github.com/spf13/cobra"
)

var (
	Logger = log.New("-")

	Config = new(config.Config)

	cfgFile string
)

func getApp() *app.App {
	db, err := database.GetDb(Config)
	if err != nil {
		log.Fatalf("unable to load db")
	}

	Queue := queue.GetConnection(Config)

	appInstance := app.App{
		Config: Config,
		Logger: Logger,
		Db:     db,
		Amqp:   Queue,
	}
	return &appInstance
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar event service",
	Long:  `Calendar event service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Use calendar [command]\nRun 'calendar --help' for usage.\n")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.example.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("root execute error: %v", err)
	}
}
