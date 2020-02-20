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
	"strings"

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
		if err != nil {
			return err
		}

		repo := cfg.Repos[i].Location
		if err := checkChanges(repo); err != nil {
			return err
		}

		msg, _ := cmd.PersistentFlags().GetString("message")
		if err := push(repo, msg); err != nil {
			return err
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

func checkChanges(repo string) error {

	out, err := common.ExecCmd("git", "-C", repo, "status", "--porcelain")
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(out)) == 0 {
		return errors.New("no changes detected: nothing to do!")
	}

	Print("Following changes have been detected for", repo)
	Print(out)
	return nil
}

func push(repo, msg string) error {
	for _, args := range [][]string{
		[]string{"add", "-A", "."},
		[]string{"commit", "-m", msg},
		[]string{"push", "origin", "master"},
	} {
		nargs := append([]string{"-C", repo}, args...)
		out, err := common.ExecCmd("git", nargs...)
		if err != nil {
			return err
		}
		Print(out)
	}

	return nil
}
