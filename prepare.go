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
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/semver"
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
	if err := os.MkdirAll(nuvDir, 0755); err != nil {
		return "", err
	}
	localDir, err := homedir.Expand("~/.nuv/olaris")
	if err != nil {
		return "", err
	}
	//fmt.Println(localDir)

	// Updating existing tools
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
		if err != nil {
			if err.Error() == "already up-to-date" {
				fmt.Println("Your tasks are already up to date!")
				return localDir, nil
			}
			return "", err
		}

		fmt.Println("Nuvfiles updated successfully")
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

	fmt.Println("Nuvfiles downloaded successfully")

	createLatestCheckFile(nuvDir)

	// clone
	return localDir, nil
}

func pullTasks(force, silent bool) error {
	// download from github
	localDir, err := downloadTasksFromGitHub(force, silent)
	debug("localDir", localDir)
	if err != nil {
		return err
	}

	// validate NuvVersion semver against nuvroot.json
	nuvRoot, err := readNuvRootFile(localDir)
	if err != nil {
		return err
	}

	// check if the version is up to date
	nuvVersion, err := semver.NewVersion(NuvVersion)
	if err != nil {
		// in development mode, we don't have a valid semver version
		warn("Unable to validate nuv version", NuvVersion, ":", err)
		return nil
	}

	nuvRootVersion, err := semver.NewVersion(nuvRoot.Version)
	if err != nil {
		warn("Unable to validate nuvroot.json version", nuvRoot.Version, ":", err)
		return nil
	}

	// check if the version is up to date, if not warn the user
	if nuvVersion.LessThan(nuvRootVersion) {
		fmt.Printf("Your nuv version (%v) is older than the required version in nuvroot.json (%v).\n", nuvVersion, nuvRootVersion)
		fmt.Println("Attempting to update nuv...")
		autoCLIUpdate()
	}

	return nil
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
		trace("found searching up:", search)
		return search, nil
	}

	// is there  olaris folder?
	olaris := joinpath(cur, "olaris")
	if exists(cur, "olaris") && exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
		trace("found sub olaris:", olaris)
		return olaris, nil
	}

	// is there an olaris folder in ~/.nuv ?
	olaris, err = homedir.Expand("~/.nuv/olaris")
	if err == nil && exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
		trace("found sub ~/.nuv/olaris:", olaris)
		return olaris, nil
	}

	// is there an olaris folder in NUV_BIN?
	nuvBin := os.Getenv("NUV_BIN")
	if nuvBin != "" {
		olaris = joinpath(nuvBin, "olaris")
		if exists(olaris, NUVFILE) && exists(olaris, NUVROOT) {
			trace("found sub NUV_BIN olaris:", olaris)
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

func autoCLIUpdate() {
	cmd := exec.Command("nuv", "update", "cli")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	// we don't care about the error, the subcommand in olaris shows it
}
