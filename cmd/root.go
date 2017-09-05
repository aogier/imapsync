// Copyright © 2017 Alessandro Ogier <alessandro.ogier@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
)

var cfgFile, host1, user1, pass1 string
var port1, pool1 int
var tls1, ssl1 bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "imapsync",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if host1 == "" {
			log.Fatal("Must provide host1, exiting ...")
		}
		if user1 == "" {
			log.Fatal("Must provide user1, exiting ...")
		}

		if ssl1 && tls1 {
			log.Fatal("ssl1 and tls1 are mutually exclusive")
		}

		if ssl1 {
			port1 = 993
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.imapsync.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&host1, "host1", "",
		"Source IMAP server")
	RootCmd.PersistentFlags().IntVar(&port1, "port1", 143,
		"Port to connect on host1. Default 143, 993 if --ssl1")
	RootCmd.PersistentFlags().StringVar(&user1, "user1", "",
		"Username on host1")
	RootCmd.PersistentFlags().StringVar(&pass1, "pass1", "",
		"Password on host1")

	RootCmd.PersistentFlags().BoolVar(&tls1, "tls1", false,
		"Force TLS connection on host1")
	RootCmd.PersistentFlags().BoolVar(&ssl1, "ssl1", false,
		"Force SSL connection on host1")

	RootCmd.PersistentFlags().IntVar(&pool1, "pool1", runtime.NumCPU(),
		"Size of connection pool against host1")

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

		// Search config in home directory with name ".imapsync" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".imapsync")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
