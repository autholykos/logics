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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all configured repository by name",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg := &Conf{}
		if err := viper.Unmarshal(cfg); err != nil {
			return err
		}

		for _, repo := range cfg.Repos {
			Print(repo.Name)
		}
		return nil
	},
}

func init() {
	downloadCmd.AddCommand(listCmd)
}
