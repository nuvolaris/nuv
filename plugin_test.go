package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
func setupPluginTest(dir string, t *testing.T) string {
	t.Helper()
	// create the olaris-test folder
	olarisTestDir := filepath.Join(dir, "olaris-test")
	if err := os.MkdirAll(olarisTestDir, 0755); err != nil {
		t.Fatal(err)
	}

	// copy the nuvroot.json from tests/olaris into the olaris-test folder
	nuvRootJSON := filepath.Join("tests", "olaris", "nuvroot.json")
	if err := copyFile(nuvRootJSON, filepath.Join(olarisTestDir, "nuvroot.json")); err != nil {
		t.Fatal(err)
	}
	// copy nuvfile.yml from tests/olaris into the olaris-test folder
	nuvfileYML := filepath.Join("tests", "olaris", "nuvfile.yml")
	if err := copyFile(nuvfileYML, filepath.Join(olarisTestDir, "nuvfile.yml")); err != nil {
		t.Fatal(err)
	}
	return olarisTestDir
}

func TestFindFolderInFolders(t *testing.T) {
	t.Run("success: plugin task found in ./olaris-test", func(t *testing.T) {
		tempDir := t.TempDir()
		plgFolder := setupPluginTest(tempDir, t)

		err := findPluginTask(tempDir, "grep")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if os.Getenv("NUV_ROOT") != plgFolder {
			t.Errorf("Expected NUV_ROOT: %s, got: %s", plgFolder, os.Getenv("NUV_ROOT"))
		}
	})

	t.Run("success: plugin task found in ~/.nuv/olaris-test", func(t *testing.T) {
		dir, _ := homedir.Expand("~/.nuv")
		// create dir/olaris-test folder
		plgFolder := setupPluginTest(dir, t)
		defer os.RemoveAll(plgFolder)

		err := findPluginTask(dir, "grep")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if os.Getenv("NUV_ROOT") != plgFolder {
			t.Errorf("Expected NUV_ROOT: %s, got: %s", plgFolder, os.Getenv("NUV_ROOT"))
		}
	})

	t.Run("error: no plugins folder found (olaris-*)", func(t *testing.T) {
		tempDir := t.TempDir()

		// Test when the folder is not found
		err := findPluginTask(tempDir, "grep")
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	t.Run("error: folder found but no plugin task found", func(t *testing.T) {
		tempDir := t.TempDir()
		plgFolder := setupPluginTest(tempDir, t)

		err := findPluginTask(tempDir, "grep-wrong")
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if os.Getenv("NUV_ROOT") == plgFolder {
			t.Errorf("Expected NUV_ROOT to not be set to: %s", plgFolder)
		}

	})
}