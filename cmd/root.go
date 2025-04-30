package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	server    string
	username  string
	password  string
	namespace string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nacos-cli",
	Short: "A CLI tool for Nacos server",
	Long: `nacos-cli is a command line interface for the Nacos server.
It provides functionality for configuration management, service discovery,
and namespace management.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nacos-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&server, "server", "http://localhost:8848", "Nacos server address")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Nacos server username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Nacos server password")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Nacos namespace ID")

	// Bind flags to viper
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nacos-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nacos-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
