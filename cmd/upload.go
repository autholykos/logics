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

	"github.com/autholykos/logics/pkg/common"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload your modification to the remote repository",
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
		fmt.Println(i)
		if err != nil {
			return err
		}

		msg, _ := cmd.PersistentFlags().GetString("message")

		for _, args := range [][]string{
			[]string{"add", "-A", "."},
			[]string{"commit", "-m", msg},
			[]string{"push", "origin", "master"},
		} {
			nargs := append([]string{"-C", cfg.Repos[i].Location}, args...)
			out, err := common.ExecCmd("git", nargs...)
			if err != nil {
				return err
			}
			fmt.Println(out)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	uploadCmd.PersistentFlags().StringP("message", "m", "committing work on Logic", "specify a message for your commit")
}
