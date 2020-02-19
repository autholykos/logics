/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type (
	Repo struct {
		Name     string `yaml:"name"`
		Location string `yaml:"location"`
	}

	Conf struct {
		SharedFolder  string `yaml:"sharedfolder"`
		ProjectFolder string `yaml:"projectfolder"`
		Repos         []Repo `yaml:"repos,flow"`
	}
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logics",
	Short: "LOGIc Control System - a version control system for distributed teams of music producers using Logic DAW",
	Long: `logics is a version control system for distributed teams of music producers.
It uses a shared (dropbox) folder as remote repository.
Internally, files are tracked through git with large file support (git-lfs) and the excellent lfs-folderstore adapter for shared folders.
For more information visit https://github.com/autholykos/logics
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// silence the annoying help on error
		cmd.SilenceUsage = true
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func WriteYaml(conf *Conf) error {
	m, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cfgFile, m, 0644); err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logics.yml)")
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

		// Search config in home directory with name ".logics" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".logics")
		cfgFile = path.Join(home, ".logics.yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}
