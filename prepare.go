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
	nuvDir, err := homedir.Expand("~/.nuv")
	if err != nil {
		return "", err
	}
	os.MkdirAll(nuvDir, 0755)
	localDir, err := homedir.Expand("~/.nuv/olaris")
	if err != nil {
		return "", err
	}
	//fmt.Println(localDir)

	// Updating an exiting tools
	// TODO: wait 24 hours...

	if exists(nuvDir, "olaris") {
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
// from github - it should contain a file nuvfile.yml and a file nuvtools.yml in the root
func locateNuvRoot(cur string) (string, error) {
	cur, err := filepath.Abs(cur)
	if err != nil {
		return "", err
	}

	// search the root from here
	search := locateNuvRootSearch(cur)
	if search != "" {
		return search, nil
	}

	// is there  olaris folder?
	olaris := joinpath(cur, "olaris")
	if exists(cur, "olaris") && exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
		return olaris, nil
	}

	// is there an olaris folder in ~/.nuv ?
	olaris, err = homedir.Expand("~/.nuv/olaris")
	if err == nil && exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
		return olaris, nil
	}

	// is there an olaris folder in NUV_BIN?
	nuvBin := os.Getenv("NUV_BIN")
	if nuvBin != "" {
		olaris = joinpath(nuvBin, "olaris")
		if exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
			return olaris, nil
		}
	}

	return "", fmt.Errorf("we cannot find nuvfiles, download them with nuv -update")
}

// locateNuvRootSearch search for `nuvfiles.yml`
// and goes up looking for a folder with also `nuvroot.json`
func locateNuvRootSearch(cur string) string {
	debug("locateNuvRootSearch:", cur)
	// exits nuvfile.yml? if not, go up until you find it
	if !exists(cur, NUVFILE) {
		return ""
	}
	if exists(cur, NUVROOT) {
		return cur
	}
	parent := parent(cur)
	if parent == "" {
		return ""
	}
	return locateNuvRootSearch(parent)
}
