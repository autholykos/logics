package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-getter"
)

var gitLFSReleaseURL = "https://github.com/git-lfs/git-lfs/releases/download/v2.10.0/git-lfs-darwin-amd64-v2.10.0.tar.gz"
var lfsFolderstoreURL = "https://github.com/sinbad/lfs-folderstore/releases/download/v1.0.0/lfs-folderstore-darwin-amd64-v1.0.0.zip"

// InstallGitLFS downloads the git-lfs package from github and installs it by
// moving it on the path (/usr/local/bin folder)
func InstallGitLFS(tmpDir string) error {
	// for some reason GitLFS package gets installed on the tmp folder
	// bypassing the package name. We work around that by adding a folder to
	// the tmpDir
	baseDir := path.Join(tmpDir, "git-lfs.pkg")
	log.WithField("tmpDir", baseDir).Debugln("downloading the git-lfs package")
	if err := Install(gitLFSReleaseURL, "", "git-lfs", baseDir); err != nil {
		return fmt.Errorf("error in installing git-lfs: %v", err)
	}

	log.Debugln("executing `git lfs install`")
	if _, err := ExecCmd("git", "lfs", "install"); err != nil {
		return fmt.Errorf("error in executing `git lfs install`: %v", err)
	}

	log.Debugln("git-lfs successfully installed")
	return nil
}

// InstallLFSFolderstore downloads the git-lfs package from github and installs it by
// moving it on the path (/usr/local/bin folder)
func InstallLFSFolderstore(tmpDir string) error {
	if err := Install(lfsFolderstoreURL, "lfs-folderstore-darwin-amd64", "lfs-folderstore", tmpDir); err != nil {
		return fmt.Errorf("error in installing lfs-folderstore: %v", err)
	}

	return nil
}

// Install downloads a package, decompress it and moves it into the path (at
// /usr/local/bin)
func Install(srcURL, pack, name, tmpDir string) error {
	log.WithFields(log.Fields{
		"tmp-dir": tmpDir,
		"pack":    pack,
		"name":    name,
	}).Debugln("downloading the package")

	// if the target dir does not exist, it gets created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return err
		}
		log.Debugf("created %s folder\n", tmpDir)
	}

	// artifact is the fully qualified path to the release
	artifact := path.Join(tmpDir, pack, name)
	log.WithField("artifact", artifact).Debugln("fully qualified name calculated")

	client := &getter.Client{
		Ctx:  context.Background(),
		Dst:  tmpDir,
		Mode: getter.ClientModeAny,
		Src:  srcURL,
	}

	if err := client.Get(); err != nil {
		return err
	}

	// checking that the file has been unpacked correctly and is available on
	// the tmp folder
	if _, err := os.Stat(artifact); os.IsNotExist(err) {
		return errors.New("download failed")
	}

	if _, err := ExecCmd("/bin/mv", artifact, "/usr/local/bin/"); err != nil {
		return fmt.Errorf("error in moving %s to /usr/local/bin: %v", artifact, err)
	}

	return nil
}
