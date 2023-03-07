package main

import (
	"os"
	"path/filepath"
	"strings"
)

// default files
const NUVFILE = "nuvfile.yml"
const NUVTOOLS = "nuvtools.json"
const NUVOPTS = "nuvopts.txt"

// repo where download tasks
const NUVREPO = "http://github.com/nuvolaris/olaris"

// branch where download tasks
// defaults to test - will be changed in compilation
var NuvBranch = "test"

// get defaults

func getNuvRoot() string {
	root := os.Getenv("NUV_ROOT")
	if root == "" {
		dir, err := os.Getwd()
		if err == nil {
			root, err = locateNuvRoot(dir, false)
		}
		if err != nil {
			panic(err)
		}
	}
	return root
}

func getNuvRepo() string {
	repo := os.Getenv("NUV_REPO")
	if repo == "" {
		repo = NUVREPO
	}
	return repo
}

func getNuvBranch() string {
	branch := os.Getenv("NUV_BRANCH")
	if branch == "" {
		branch = NuvBranch
	}
	return branch
}

// utils
func join(dir string, file string) string {
	return filepath.Join(dir, file)
}

func split(s string) []string {
	return strings.Fields(s)
}

func exists(dir string, file string) bool {
	_, err := os.Stat(join(dir, file))
	return !os.IsNotExist(err)
}

func parent(dir string) string {
	return filepath.Dir(dir)
}

func readfile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(dat)
}
