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
	"os"

	"github.com/spf13/cobra"
)

var target string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a Logic project hosted on a shared folder",
	Long: `Install a Logic project in the default Logic folder (or in a the target folder if that is specified). For example:

  logics install capelli-curti				# install project "capelli-curti" on defaul Logic directory
  logics install -t /path/to/folder capelli-curti	# install project "capelli-curti" on /path/to/folder
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.
	// TODO: default should be looked up on the .yaml config file written
	// after the setup phase
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	defaultDir := fmt.Sprintf("%s/Music/Logic", home)

	installCmd.PersistentFlags().StringVarP(&target, "target-folder", "t", defaultDir, "specify a different target folder than the default one")
}
