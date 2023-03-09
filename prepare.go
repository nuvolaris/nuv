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

func downloadTasksFromGitHub(force bool, silent bool) (string, error) {
	repoURL := getNuvRepo()
	branch := getNuvBranch()
	localDir, err := homedir.Expand("~/.nuv/olaris")
	//fmt.Println(localDir)
	if err != nil {
		return "", err
	}

	// Updating an exiting tools
	// TODO: wait 24 hours...

	if exists(localDir, NUVTOOLS) {
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

	// Clone the repo if not existing
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

// locateNuvRoot locate the folder where starts execution
// it can be a parent folder of the current folder or it can be downloaded
// from github - it should contain a file Nuvfile and a file Nuvtools
func locateNuvRoot(cur string, inHomeDir bool) (string, error) {
	cur, err := filepath.Abs(cur)
	if err != nil {
		return "", err
	}
	// exits nuvtile.yml?
	if !exists(cur, NUVFILE) {
		if exists(cur, "olaris") {
			return locateNuvRoot(join(cur, "olaris"), inHomeDir)
		}
		if inHomeDir {
			return "", fmt.Errorf("cannot find nuv root dir and cannot download it from github")
		}
		dir, err := downloadTasksFromGitHub(true, true)
		if err != nil {
			return "", err
		}
		return locateNuvRoot(dir, true)
	}

	// found nuvtoools.json? if not, go up
	if !exists(cur, NUVTOOLS) {
		return locateNuvRoot(parent(cur), inHomeDir)
	}
	return cur, nil
}
