// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mitchellh/go-homedir"
)

func join(dir string, file string) string {
	return filepath.Join(dir, file)
}

func exists(dir string, file string) bool {
	_, err := os.Stat(join(dir, file))
	return !os.IsNotExist(err)
}

func parent(dir string) string {
	return filepath.Dir(dir)
}

func loadTools() {
	//return make(map[string]string)
}

func downloadTasksFromGitHub(force bool, silent bool) (string, error) {
	repoURL := os.Getenv("NUV_OLARIS_REPO")
	if repoURL == "" {
		repoURL = NuvOlarisRepo
	}
	localDir, err := homedir.Expand("~/.nuv/olaris")
	//fmt.Println(localDir)
	if err != nil {
		return "", err
	}

	branch := NuvOlarisBranch
	if branch == "" {
		branch = os.Getenv("NUV_OLARIS_BRANCH")
	}

	if exists(localDir, "Nuvtools") {
		fmt.Println("Updating tasks...")
		r, err := git.PlainOpen(localDir)
		if err != nil {
			return "", err
		}
		// Get the working directory for the repository
		w, err := r.Worktree()
		if err != nil {
			return "", err
		}

		// Pull the latest changes from the origin remote and merge into the current branch
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err.Error() == "already up-to-date" {
			return localDir, nil
		}
		if err != nil {
			return "", err
		}
		return localDir, nil
	}

	//fmt.Println(repoURL, branch)
	ref := plumbing.NewBranchReferenceName(branch)
	//fmt.Println(ref)
	cloneOpts := &git.CloneOptions{
		URL:           repoURL,
		Progress:      os.Stderr,
		ReferenceName: ref, // Specify the branch to clone
	}

	fmt.Println("Cloning tasks...")
	_, err = git.PlainClone(localDir, false, cloneOpts)
	if err != nil {
		return "", err
	}
	// clone
	return localDir, nil
}

// prepareTaskFolderAndTools locate the folder where starts execution
// it can be a parent folder of the current folder or it can be downloaded
// from github - it should contain a file Nuvfile and a file Nuvtools
func prepareTaskFolderAndTools(cur string, inHomeDir bool) (string, error) {
	cur, err := filepath.Abs(cur)
	if err != nil {
		return "", err
	}
	if !exists(cur, "Nuvfile") {
		if exists(cur, "olaris") {
			return prepareTaskFolderAndTools(join(cur, "olaris"), inHomeDir)
		}
		if inHomeDir {
			return "", fmt.Errorf("cannot find nuv root dir and cannot download it from github")
		}
		dir, err := downloadTasksFromGitHub(true, true)
		if err != nil {
			return "", err
		}
		return prepareTaskFolderAndTools(dir, true)
	}

	// exists
	if !exists(cur, "Nuvtools") {
		return prepareTaskFolderAndTools(parent(cur), inHomeDir)
	}

	loadTools()
	return cur, nil
}
