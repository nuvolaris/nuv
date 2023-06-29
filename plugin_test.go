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
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
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

func TestFindPluginTask(t *testing.T) {
	t.Run("success: plugin task found in ./olaris-test", func(t *testing.T) {
		tempDir := t.TempDir()
		plgFolder := setupPluginTest(tempDir, t)

		fld, err := findTaskInPlugins(tempDir, "grep")
		require.NoError(t, err)
		require.Equal(t, plgFolder, fld)

		// if os.Getenv("NUV_ROOT") != plgFolder {
		// 	t.Errorf("Expected NUV_ROOT: %s, got: %s", plgFolder, os.Getenv("NUV_ROOT"))
		// }
	})

	t.Run("success: plugin task found in ~/.nuv/olaris-test", func(t *testing.T) {
		dir, _ := homedir.Expand("~/.nuv")
		// create dir/olaris-test folder
		plgFolder := setupPluginTest(dir, t)
		defer os.RemoveAll(plgFolder)

		fld, err := findTaskInPlugins(dir, "grep")
		require.NoError(t, err)
		require.Equal(t, plgFolder, fld)

		// if os.Getenv("NUV_ROOT") != plgFolder {
		// 	t.Errorf("Expected NUV_ROOT: %s, got: %s", plgFolder, os.Getenv("NUV_ROOT"))
		// }
	})

	t.Run("error: no plugins folder found (olaris-*)", func(t *testing.T) {
		tempDir := t.TempDir()

		// Test when the folder is not found
		fld, err := findTaskInPlugins(tempDir, "grep")
		require.Error(t, err)
		require.Empty(t, fld)
	})

	t.Run("error: existing plugin folder but no plugin task found", func(t *testing.T) {
		tempDir := t.TempDir()
		_ = setupPluginTest(tempDir, t)

		fld, err := findTaskInPlugins(tempDir, "grep-wrong")
		require.Error(t, err)
		require.Empty(t, fld)
	})
}

func TestNewPlugins(t *testing.T) {
	t.Run("create plugins struct with valid local dir", func(t *testing.T) {
		tempDir := t.TempDir()
		plgFolder := setupPluginTest(tempDir, t)

		p, err := newPlugins(tempDir)
		require.NoError(t, err)
		require.NotNil(t, p)
		require.Len(t, p.local, 1)
		require.Equal(t, plgFolder, p.local[0])
	})

	t.Run("non existent local dir results in empty local field", func(t *testing.T) {
		localDir := "/path/to/nonexistent/dir"
		p, err := newPlugins(localDir)
		require.NoError(t, err)
		require.NotNil(t, p)
		require.Len(t, p.local, 0)
	})
}

func Example_pluginsPrint() {
	p := plugins{
		local: make([]string, 0),
		nuv:   make([]string, 0),
	}
	p.print()
	// Output
	// No plugins installed. Use 'nuv -plugin' to add new ones.
}

func TestCheckGitRepo(t *testing.T) {
	tests := []struct {
		url          string
		expectedRepo bool
		expectedName string
	}{
		{
			url:          "https://github.com/giusdp/olaris-test",
			expectedRepo: true,
			expectedName: "olaris-test",
		},
		{
			url:          "https://github.com/giusdp/olaris-test.git",
			expectedRepo: true,
			expectedName: "olaris-test",
		},
		{
			url:          "git@github.com:giusdp/olaris-test.git",
			expectedRepo: true,
			expectedName: "olaris-test",
		},
		{
			url:          "https://github.com/giusdp/some-repo",
			expectedRepo: false,
			expectedName: "",
		},
		{
			url:          "https://github.com/giusdp/olaris-repo.git",
			expectedRepo: true,
			expectedName: "olaris-repo",
		},
		{
			url:          "https://github.com/olaris-1234/repo",
			expectedRepo: false,
			expectedName: "",
		},
		{
			url:          "https://github.com/giusdp/another-repo.git",
			expectedRepo: false,
			expectedName: "",
		},
	}

	for _, test := range tests {
		isOlarisRepo, repoName := checkGitRepo(test.url)
		require.Equal(t, test.expectedRepo, isOlarisRepo)
		require.Equal(t, test.expectedName, repoName)
	}
}
