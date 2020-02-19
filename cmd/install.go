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
	"path/filepath"
	"strings"

	"github.com/autholykos/logics/pkg/common"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringP("projectfolder", "p", "", "specify a target project folder")
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a Logic project hosted on a shared folder",
	Long: `Install a Logic project in the default Logic folder (or in a the target folder if that is specified). For example:

  logics install # install checks for projects within the shared folder and install it on the default Logic directory
  logics install -p /path/to/folder # install project "capelli-curti" on /path/to/folder
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkSetup(); os.IsNotExist(err) {
			return errors.New("No config file found for logics. Please run `logics setup` first")
		}

		folder, _ := cmd.PersistentFlags().GetString("projectfolder")
		if strings.TrimSpace(folder) == "" {
			folder = viper.GetString("projectfolder")
		}

		if _, err := os.Stat(folder); os.IsNotExist(err) {
			return fmt.Errorf("project folder %s does not exist. Please (re)run `logics setup` or specify a different folder", folder)
		}

		viper.Set("targetdir", folder)
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		sharedDir := viper.GetString("sharedfolder")
		cfg := &Conf{}
		if err := viper.Unmarshal(cfg); err != nil {
			return err
		}
		remoteRepo, err := selectRepo(sharedDir, cfg)
		if err != nil {
			return err
		}
		localDir := viper.GetString("targetdir")
		basename := strings.TrimSuffix(filepath.Base(remoteRepo), ".git")
		localRepo := path.Join(localDir, basename)
		localRepoGit := fmt.Sprintf("%s.git", localRepo)

		if err := cloneRepo(localRepo, remoteRepo); err != nil {
			return err
		}

		if err := configureLFSFolderstore(localRepo, remoteRepo); err != nil {
			return err
		}

		fmt.Println("new repository installed and configured")
		cfg.Repos = append(cfg.Repos, Repo{
			Name:     basename,
			Location: localRepoGit,
		})

		yfg, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(viper.ConfigFileUsed(), yfg, 0644); err != nil {
			return err
		}

		fmt.Println("preferences saved")
		return nil
	},
}

func cloneRepo(localRepo, remoteRepo string) error {
	//out, err := common.ExecCmd("git", "-C", localDir, "clone", remoteRepo)
	out, err := common.ExecCmd("git", "clone", remoteRepo, localRepo)
	if err != nil {
		return fmt.Errorf("error in cloning the repo: %v", err)
	}
	fmt.Println(out)
	return nil
}

func configureLFSFolderstore(localRepo, remoteRepo string) error {
	if err := execGit(localRepo, "config", "--add", "lfs.customtransfer.lfs-folder.path", "lfs-folderstore"); err != nil {
		return err
	}
	if err := execGit(localRepo, "config", "--add", "lfs.customtransfer.lfs-folder.args", remoteRepo); err != nil {
		return err
	}
	if err := execGit(localRepo, "config", "--add", "lfs.standalonetransferagent", "lfs-folder"); err != nil {
		return err
	}
	if err := execGit(localRepo, "reset", "--hard", "master"); err != nil {
		return err
	}
	fmt.Println("lfs-folderstore configured")
	return nil
}

func execGit(localRepo string, args ...string) error {
	args = append([]string{"-C", localRepo}, args...)
	out, err := common.ExecCmd("git", args...)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func selectRepo(sharedDir string, conf *Conf) (string, error) {
	projects := make([]string, 0)
	files, err := ioutil.ReadDir(sharedDir)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		// we add all base names that have a .git folder in it
		if _, err := os.Stat(path.Join(sharedDir, f.Name(), ".git")); !os.IsNotExist(err) {
			repo := path.Join(sharedDir, f.Name())
			if isAlreadyCloned(repo, conf) {
				continue
			}
			// NOTE: projects are in the form [/path/to/project.git]
			projects = append(projects, path.Join(sharedDir, f.Name()))
		}
	}

	if len(projects) == 0 {
		return "", fmt.Errorf("no new project found in %s", sharedDir)
	}

	prompt := promptui.Select{
		Label: "select which project you want to install",
		Items: projects,
	}

	_, repo, err := prompt.Run()
	return repo, err
}

func isAlreadyCloned(repo string, conf *Conf) bool {
	if conf.Repos == nil {
		return false
	}
	for _, r := range conf.Repos {
		if filepath.Base(r.Location) == filepath.Base(repo) {
			return true
		}
	}
	return false
}

/* this is for setting up a project
func track(localRepo string) {
	audioTypes := []string{"*.wav", "*aif", "*.mp3"}
	for _, audioType := range audioTypes {
		out, err := common.ExecCmd("git", "-C", localRepo, "lfs", "track", audioType)
		if err != nil {
			return err
		}
		fmt.Println(out)
	}
		out, err := common.ExecCmd("git", "-C", localRepo, "add", ".gitattributes")
		if err != nil {
			return err
		}
		fmt.Println(out)
		out, err := common.ExecCmd("git", "-C", localRepo, "commit", "-m", "add wav, aif and mp3 to the LFS")
		if err != nil {
			return err
		}
		fmt.Println(out)
}
*/

func checkSetup() error {
	hd, _ := os.UserHomeDir()
	_, err := os.Stat(path.Join(hd, ".logics.yml"))
	return err
}
