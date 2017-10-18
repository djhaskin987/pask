// Copyright © 2017 Daniel Jay Haskin <djhaskin987@gmail.com>
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
	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pask",
	Short: "PAcKaged tASKs",
	Long: `Install files and run tasks. Make your build
less painful and more fun!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Printf("%v\n", viper.Get("spec"))
	//	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	viper.SetEnvPrefix("PASK")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config", "c", "", "config file (default is $HOME/.pask)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	RootCmd.PersistentFlags().StringP("base", "p", "./", "Base project path")
	viper.BindPFlag("base", RootCmd.PersistentFlags().Lookup("base"))
	var baseDefault string
	if pwd, err := os.Getwd(); err != nil {
		baseDefault = "."
	} else {
		baseDefault = pwd
	}
	viper.SetDefault("base", baseDefault)

	RootCmd.PersistentFlags().StringP("spec", "s", "", "Pask spec file")
	viper.BindPFlag("spec", RootCmd.PersistentFlags().Lookup("spec"))
	viper.SetDefault("spec", path.Join(baseDefault,
		"pask",
		"spec.hcl"))
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
			log.Fatalln(err)
		}

		// Search config in home directory with name ".pask" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pask")
		viper.SetConfigType("toml")

		// In the succeeding line, the reader will note that
		// ReadInConfig returns an error, and I am intentionally
		// ignoring it.
		//
		// If the config file couldn't be read, *I don't care*. But
		// if it can be, it should be.
		err = viper.ReadInConfig()
		if err != nil {
			log.Println("Error reading config file: ", err)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			viper.ReadInConfig()
		})
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
