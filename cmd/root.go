// arkadiko
// https://github.com/topfreegames/arkadiko
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile is the configuration file used for running a command
var ConfigFile string

// Verbose determines how verbose arkadiko will run under
var Verbose int

// RootCmd is the root command for arkadiko CLI application
var RootCmd = &cobra.Command{
	Use:   "arkadiko",
	Short: "arkadiko bridges http to mqtt",
	Long:  `arkadiko bridges http to mqtt.`,
}

// Execute runs RootCmd to initialize arkadiko CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().IntVarP(
		&Verbose, "verbose", "v", 0,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)

	RootCmd.PersistentFlags().StringVarP(
		&ConfigFile, "config", "c", "./config/local.yml",
		"config file (default is ./config/local.yml",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(ConfigFile)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("arkadiko")
	viper.SetConfigName(".arkadiko") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
