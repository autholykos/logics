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
	"github.com/autholykos/logics/pkg/common"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "update your local repository with all changes performed remotely",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &Conf{}
		if err := viper.Unmarshal(cfg); err != nil {
			return err
		}

		projects := make([]string, 0)
		for _, repo := range cfg.Repos {
			projects = append(projects, repo.Name)
		}

		prompt := promptui.Select{
			Label: "select which project you want to sync",
			Items: projects,
		}

		i, _, err := prompt.Run()
		if err != nil {
			return err
		}

		out, err := common.ExecCmd("git", "-C", cfg.Repos[i].Location, "pull", "origin", "master")
		if err != nil {
			return err
		}

		Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
