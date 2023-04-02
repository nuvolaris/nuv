// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
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
	"time"

	"github.com/go-git/go-git/v5"
)

const LATESTCHECK = ".latestcheck"

func checkUpdated(base string, timeInterval time.Duration) {
	latest_check_path := joinpath(base, LATESTCHECK)
	olaris_path := joinpath(base, "olaris")

	// if no olaris dir, no update check
	if exists(base, "olaris") && !isDir(olaris_path) {
		return
	}

	// get info on latest_check file
	file, err := os.Stat(latest_check_path)
	if err != nil {
		// if latest_check file doesn't exist, or error accessing it, create it
		createLatestCheckFile(base)
		return
	}

	mtime := file.ModTime()
	now := time.Now().Local()
	diff := now.Sub(mtime)

	if diff >= timeInterval {
		fmt.Println("Checking for updates...")
		// touch latest_check file ONLY if enough time has passed
		touchLatestCheckFile(latest_check_path)

		// check if remote olaris is newer
		if checkRemoteOlarisNewer(olaris_path) {
			fmt.Print("New tasks available! Use 'nuv -update' to update.\n\n")
		} else {
			fmt.Print("Tasks up to date!\n\n")
		}
	}
}

func createLatestCheckFile(base string) {
	// create latest_check file
	_, err := os.Create(joinpath(base, LATESTCHECK))
	if err != nil {
		warn("failed to set latest_check file", err)
	}
}

func touchLatestCheckFile(latest_check_path string) {
	currentTime := time.Now().Local()
	err := os.Chtimes(latest_check_path, currentTime, currentTime)
	if err != nil {
		warn("failed to set latest update check", err)
	}
}

func checkRemoteOlarisNewer(olaris_path string) bool {
	isRemoteNewer := false

	repo, err := git.PlainOpen(olaris_path)
	if err != nil {
		warn("failed to check olaris folder", err)
		return isRemoteNewer
	}

	localRef, _ := repo.Head()

	remote, _ := repo.Remote("origin")
	remote.Fetch(&git.FetchOptions{})
	remoteRefs, _ := remote.List(&git.ListOptions{})

	// check ref is in refs
	for _, remoteRef := range remoteRefs {
		if localRef.Name().String() == remoteRef.Name().String() {
			// is hash different?
			if localRef.Hash().String() != remoteRef.Hash().String() {
				isRemoteNewer = true
			}
		}
	}
	return isRemoteNewer
}