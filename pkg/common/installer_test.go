package common_test

import (
	"os"
	"testing"

	"github.com/autholykos/logics/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestGitLFSInstall(t *testing.T) {
	defer os.RemoveAll("/tmp/logics")

	if !assert.NoError(t, os.MkdirAll("/tmp/logics", 0755)) {
		t.FailNow()
	}
	if !assert.NoError(t, common.InstallGitLFS("/tmp/logics")) {
		t.FailNow()
	}

	_, err := os.Stat("/usr/local/bin/git-lfs")
	assert.NoError(t, err)
}

func TestLFSFolderstoreInstall(t *testing.T) {
	defer os.RemoveAll("/tmp/logics")

	if !assert.NoError(t, os.MkdirAll("/tmp/logics", 0755)) {
		t.FailNow()
	}
	if !assert.NoError(t, common.InstallLFSFolderstore("/tmp/logics")) {

		t.FailNow()
	}

	_, err := os.Stat("/usr/local/bin/lfs-folderstore")
	assert.NoError(t, err)
}
