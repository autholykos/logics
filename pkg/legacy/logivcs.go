package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	getter "github.com/hashicorp/go-getter"
)

const defaultFailedCode = -1

func main() {
	var dropboxDir string
	hd, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-4)
	}

	log.Println(fmt.Sprintf("running under user's home %s. Target dir %s/Music/Logic", hd, hd))
	if exitCode := execCmd("git", "--version"); exitCode != 0 {
		os.Exit(exitCode)
	}

	log.Println("is the exfalsoproject dropbox available?")
	dropboxDir = fmt.Sprintf("%s/Dropbox/logic", hd)
	if _, err := os.Stat(dropboxDir); os.IsNotExist(err) {
		log.Println("please install dropbox shared folder on this machine. User: `exfalsoproject` pwd: `...`")
		os.Exit(-2)
	}
	log.Println("exfalso dropbox folder found!")

	log.Println("installing git-lfs")
	if exitCode := execCmd("brew", "install", "git-lfs"); exitCode == defaultFailedCode {
		if ec := execCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""); ec != 0 {
			log.Fatalln("could not install homebrew. Try to install it manually by following the instruction at https://brew.sh")
			os.Exit(-2)
		}
	} else if exitCode != 0 {
		os.Exit(exitCode)
	}

	log.Println("installing lfs-folderstore")
	if err := installLFSFolderstore("https://github.com/sinbad/lfs-folderstore/releases/download/v1.0.0/lfs-folderstore-darwin-amd64-v1.0.0.zip"); err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	log.Println("cloning the Logic repo")
	remoteRepo := fmt.Sprintf("%s/capelli-curti.git", dropboxDir)
	localRepo := fmt.Sprintf("%s/Dev/capelli-curti", hd)

	execCmd("git", "-C", localRepo, "lfs", "track", "*.wav")
	execCmd("git", "-C", localRepo, "lfs", "track", "*.aif")
	execCmd("git", "-C", localRepo, "lfs", "track", "*.mp3")
	execCmd("git", "-C", localRepo, "add", ".gitattributes")
	execCmd("git", "-C", localRepo, "commit", "-m", "add wav, aif and mp3 to the LFS")

	execCmd("git", "clone", remoteRepo, localRepo)
	execCmd("git", "-C", localRepo, "config", "--add", "lfs.customtransfer.lfs-folder.path", "lfs-folderstore")
	execCmd("git", "-C", localRepo, "config", "--add", "lfs.customtransfer.lfs-folder.args", remoteRepo)
	execCmd("git", "-C", localRepo, "config", "--add", "lfs.standalonetransferagent", "lfs-folder")
	execCmd("git", "-C", localRepo, "reset", "--hard", "master")

	os.Exit(0)
}

func execCmd(name string, args ...string) int {

	// first we check if there is a git installation already
	stdout, stderr, err := runcmd(name, args...)
	if err != nil {
		exitCode := extractExitCode(err)
		if exitCode == defaultFailedCode {
			log.Fatalf("command `%s` not found. Please install it.", name)
		} else if len(stderr) > 0 {
			log.Fatalf("command failed with following message: %s", string(stderr))
		} else {
			log.Fatalf("error detected: %v", err)
		}
		return exitCode
	}

	log.Println(string(stdout))
	return 0
}

func extractExitCode(err error) int {

	// try to get the exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		return ws.ExitStatus()
	}

	// This will happen (in OSX) if `name` is not available in $PATH,
	// in this situation, exit code could not be get, and stderr will be
	// empty string very likely, so we use the default fail code, and format err
	// to string and set to stderr
	return defaultFailedCode
}

func runcmd(name string, args ...string) ([]byte, []byte, error) {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString(" ")
	for i, arg := range args {
		if i != 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(arg)
	}

	log.Println(fmt.Sprintf("running cmd: `%s`", sb.String()))
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	if err := cmd.Run(); err != nil {
		// in case of error, we return the Stderr + the error
		return outbuf.Bytes(), errbuf.Bytes(), err
	}

	// success, exitCode should be 0 if go is ok
	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func installLFSFolderstore(target string) error {
	tgd := "/tmp/logivcs/lfs-folderstore.pkg"
	artifact := tgd + "/lfs-folderstore-darwin-amd64/lfs-folderstore"
	//bin := "/usr/local/bin/lfs-folderstore"
	bin := "/usr/local/bin/lfs-folderstore"

	if _, err := os.Stat(bin); !os.IsNotExist(err) {
		log.Println("lfs-folderstore already installed in /usr/local/bin")
		return nil
	}

	defer os.RemoveAll("/tmp/logivcs/lfs-folderstore.pkg")

	client := &getter.Client{
		Ctx:  context.Background(),
		Dst:  tgd,
		Mode: getter.ClientModeAny,
		Src:  target,
	}

	if err := client.Get(); err != nil {
		return err
	}

	if _, err := os.Stat(tgd); os.IsNotExist(err) {
		//TODO: cleanup
		log.Fatalln("something went wrong with go-getter")
		os.Exit(-3)
	}

	if ec := execCmd("/bin/mv", artifact, bin); ec != 0 {
		log.Fatalln("could not copy", artifact, "in", bin)
		os.Exit(-2)
	}

	return nil
}
