# Logic Version Control Installer

This program helps to install the necessary tools for starting up and maintaining a Logic project fully available to a distributed team of musicians, sound engineers and composers. Its focus is simplicity and avoiding collaborators within a Logic project to having to toss around audio tracks and other deliverables.

## Prerequisite

* MacOSX: from Mojave up
* git (minimum version 2.0): version control system used by this program
* shared folder: the project relies on shared folder (Dropbox)  as "remote" storage where the git repositories are stored, which needs to be made available to the program in order to function properly.

## Dependencies

The program will download and install the following programs on the local system:
* git-lfs: git support for large file system
* lfs-folderstore: shared folder agent for git-lfs

## Installation

if you wish to compile the source code, you can clone this repository and then run `go install` directly in the project root. Releases will be available soon

## Usage

### Setup

The `setup` command should be called only once in order to configure the version control system. It downloads and install `git-lfs` and `lfs-folderstore`, let you specify the shared folder where the _remote_ repository is found, and the target directory where your (Logic) projects should be installed. The configuration is written on `$HOME/.logics.ylm`

```
$ logics setup`
```

### Install

The `install` command scans the shared folder for repositories not yet installed, let you select the repository you want to pull and configures `git-lfs` to track audio files. If no target directory is specified the default folder specified during setup gets used

### Download

The `download` command let you sync up with the upstream by selecting a repo and download the latest changes

```
$ logics download
```

### Upload

The `upload` command let you upload your changes to upstream if any

```
$ logics upload
```
