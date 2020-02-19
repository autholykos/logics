/*
Copyright © 2020 Emanuele Francioni <emanuele@dusk.network>

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

	"github.com/autholykos/logics/pkg/common"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dropboxFolder, logicFolder, tmpDir string

func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dropboxFolder = path.Join(hd, "Dropbox", "logic")
	logicFolder = path.Join(hd, "Music", "Logic")
	tmpDir, _ = ioutil.TempDir("", "logics")

	rootCmd.AddCommand(setupCmd)
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup logics configuration",
	Long:  `Setup the yaml file used to persist configuration attributes for using logics`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := ExecCmd("git", "--version")
		if err != nil {
			return errors.New("no git installation found")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Setup()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if err := os.RemoveAll(tmpDir); err != nil {
			return fmt.Errorf("WARNING: could not remove %s", tmpDir)
		}
		return nil
	},
}

func Setup() error {
	cfg := strings.TrimSpace(viper.ConfigFileUsed())
	if len(cfg) > 0 {
		if !common.YNPrompt(fmt.Sprintf("A setup was likely already run (and created the configuration at %s). Do you want to re-run the setup?", cfg)) {
			fmt.Println("Okidokey")
			return nil
		}
	}

	conf := &Conf{}
	sharedDir, err := setupSharedDir()
	if err != nil {
		return err
	}
	conf.SharedFolder = sharedDir

	projDir, err := setupProjectDir()
	if err != nil {
		return err
	}
	conf.ProjectFolder = projDir

	if err := WriteYaml(conf); err != nil {
		return fmt.Errorf("Something went wrong with writing config file %s, %v", cfg, err)
	}
	fmt.Println("Preferences saved", cfg)

	if err := common.InstallGitLFS(tmpDir); err != nil {
		return err
	}
	fmt.Println("git-lfs installed", cfg)

	if err := common.InstallLFSFolderstore(tmpDir); err != nil {
		return err
	}
	fmt.Println("lfs-folderstore installed", cfg)
	return nil
}

// setupSharedDir sets up the shared repository
func setupSharedDir() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Please input the shared folder path",
		Default: dropboxFolder,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	result = strings.TrimSpace(result)

	if result == "" {
		result = dropboxFolder
	}

	if err := validateDir(result); err != nil {
		return "", err
	}

	return result, nil
}

// setupProjectDir sets up the folder with the Logic projects
func setupProjectDir() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Please input your project folder",
		Default: logicFolder,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	result = strings.TrimSpace(result)

	if err := validateDir(result); err != nil {
		if !common.YNPrompt(fmt.Sprintf("Cannot find %s. Do you want to create it?", result)) {
			return "", errors.New("Cannot setup without a project folder")
		}

		if err := upsertDir(result); err != nil {
			return "", err
		}
	}

	return result, nil
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
