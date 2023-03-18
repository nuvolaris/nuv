package main

import (
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
)

func changeLatestCheckTime(base string, t time.Duration) {
	latestCheckPath := joinpath(base, ".latestcheck")
	file, err := os.Stat(latestCheckPath)
	if err != nil {
		pr("failed to get latest_check file info", err)
	}
	mtime := file.ModTime()
	mtime = mtime.Add(t)
	err = os.Chtimes(latestCheckPath, mtime, mtime)
	if err != nil {
		pr("failed to set latest_check file mtime", err)
	}
}

func resetOneCommit(repo *git.Repository) {
	commIter, _ := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})

	commIter.Next()
	secondLastCommit, _ := commIter.Next()

	w, _ := repo.Worktree()
	w.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: secondLastCommit.Hash,
	})
}

func Example_checkUpdated_uptodate() {
	// clone olaris folder into a temp folder
	tmpDir, err := os.MkdirTemp("", "nuv-test")
	if err != nil {
		pr("failed to create temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	olarisTmpPath := joinpath(tmpDir, "olaris")

	_, err = git.PlainClone(olarisTmpPath, false, &git.CloneOptions{
		URL:      getNuvRepo(),
		Progress: os.Stderr},
	)

	// run checkUpdated and check if it creates the latest_check file
	checkUpdated(tmpDir, 1*time.Second)

	if exists(tmpDir, ".latestcheck") {
		pr("latest_check file created")
	}

	// change latest_check file mtime to 2 seconds ago
	changeLatestCheckTime(tmpDir, -2*time.Second)

	// re-run checkUpdated and check output "Tasks up to date!"
	checkUpdated(tmpDir, 1*time.Second)

	// Output:
	// latest_check file created
	// Checking for updates...
	// Tasks up to date!
}

func Example_checkUpdated_outdated() {
	// clone olaris folder into a temp folder
	tmpDir, err := os.MkdirTemp("", "nuv-test")
	if err != nil {
		pr("failed to create temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	olarisTmpPath := joinpath(tmpDir, "olaris")

	repo, err := git.PlainClone(olarisTmpPath, false, &git.CloneOptions{
		URL:      getNuvRepo(),
		Progress: os.Stderr},
	)

	// run checkUpdated and check if it creates the latest_check file
	checkUpdated(tmpDir, 1*time.Second)

	if exists(tmpDir, ".latestcheck") {
		pr("latest_check file created")
	}

	// change latest_check file mtime to 2 seconds ago
	changeLatestCheckTime(tmpDir, -2*time.Second)

	// git reset olaris to a previous commit
	resetOneCommit(repo)

	// re-run checkUpdated and check output
	checkUpdated(tmpDir, 1*time.Second)

	// Output:
	// latest_check file created
	// Checking for updates...
	// New tasks available! Use 'nuv -update' to update.
}
