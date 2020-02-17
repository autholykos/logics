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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var dropboxFolder, logicFolder, configFile string

type (
	Conf struct {
		//TODO: should we support URLs based git-lfs?
		SharedFolderPath  string `yaml:"shared-folder"`
		ProjectFolderPath string `yaml:"project-folder"`
	}
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup logics configuration",
	Long:  `Setup the yaml file used to persist configuration attributes for using logics`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var result string

		cfg := strings.TrimSpace(viper.ConfigFileUsed())
		if len(cfg) > 0 {
			prompt := promptui.Select{
				Label: fmt.Sprintf("A setup was likely already run (and created the configuration at %s). Do you want to re-run the setup?", cfg),
				Items: []string{"Nay", "Yay"},
			}
			_, yayOrNay, e := prompt.Run()
			if e != nil {
				fmt.Println(e)
				return
			}

			if yayOrNay != "Yay" {
				fmt.Println("Okidokey")
				return
			}
		}

		conf := &Conf{}
		prompt := promptui.Prompt{
			Label:   "Please input the shared folder path",
			Default: dropboxFolder,
		}

		result, err = prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		result = strings.TrimSpace(result)

		if result == "" {
			result = dropboxFolder
		}

		if err := validateDir(result); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Shared folder set to", result)
		conf.SharedFolderPath = result

		prompt = promptui.Prompt{
			Label:   "Please input your project folder",
			Default: logicFolder,
		}

		result, err = prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		result = strings.TrimSpace(result)

		if err := validateDir(result); err != nil {
			sprompt := promptui.Select{
				Label: fmt.Sprintf("Cannot find %s. Do you want to create it?", result),
				Items: []string{"Yay", "Nay"},
			}
			_, yayOrNay, e := sprompt.Run()
			if e != nil {
				fmt.Println(e)
				return
			}

			if yayOrNay != "Yay" {
				fmt.Println(errors.New("Cannot setup without a project folder"))
				return
			}

			if err := upsertDir(result); err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Println("Project folder set to", result)
		conf.ProjectFolderPath = result
		m, e := yaml.Marshal(conf)
		if e != nil {
			fmt.Println(e)
			return
		}

		if err := ioutil.WriteFile(configFile, m, 0644); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Wrote preferences to", configFile)
	},
}

func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dropboxFolder = path.Join(hd, "Dropbox", "logic")
	logicFolder = path.Join(hd, "Music", "Logic")
	configFile = path.Join(hd, ".logics.yml")

	rootCmd.AddCommand(setupCmd)
}

func validateDir(dpath string) error {
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		return fmt.Errorf("could not find %s on this machine. Please make sure you use a fully qualified name (e.g. /Users/pippo/Dropbox)", dpath)
	}

	return nil
}

func upsertDir(dpath string) error {
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		if err := os.MkdirAll(dpath, 0755); err != nil {
			return err
		}
	}
	return nil
}
