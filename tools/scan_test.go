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

package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// func Test_buildActionPlan(t *testing.T) {
// 	cmd := []string{"nuv", "-js", "script.js"}

// 	t.Run("returns error if no actions folder", func(t *testing.T) {
// 		fakeFs := fstest.MapFS{"a": {Data: []byte{}}}
// 		_, err := buildActionPlan(fakeFs, cmd)
// 		if err == nil {
// 			t.Error("expected error, got nil")
// 		}
// 	})

// 	t.Run("returns one arg cmd plan if empty folder", func(t *testing.T) {
// 		fakeFs := fstest.MapFS{actionsFolder: {
// 			Mode: fs.ModeDir,
// 		}}
// 		plan, err := buildActionPlan(fakeFs, cmd)
// 		if err != nil {
// 			t.Errorf("expected nil, got %v", err)
// 		}
// 		if len(plan.args) != 1 {
// 			t.Errorf("expected 1, got %d", len(plan.args))
// 		}
// 		if plan.args[0][0] != "/actions" {
// 			t.Errorf("expected /actions, got %s", plan.args[0][0])
// 		}
// 	})

// }

func Test_checkActionsFolder(t *testing.T) {
	t.Run("returns error if actions folder does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := checkActionsFolder(tmpDir)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error if actions folder is not a folder", func(t *testing.T) {
		tmpDir := t.TempDir()
		actionsFile := filepath.Join(tmpDir, actionsFolder)
		_, err := os.Create(actionsFile)
		if err != nil {
			t.Fatalf("failed to create file %s: %v", actionsFile, err)
		}
		err = checkActionsFolder(tmpDir)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns nil if actions folder exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		actionsFile := filepath.Join(tmpDir, actionsFolder)
		err := os.Mkdir(actionsFile, 0755)
		if err != nil {
			t.Fatalf("failed to create file %s: %v", actionsFile, err)
		}

		err = checkActionsFolder(tmpDir)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func Test_getAllDirs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	expectedDirs := []string{tempDir}

	// Create some subdirectories inside the temporary directory
	subDirs := []string{"dir1", "dir2", "dir3"}
	subSubDirs := []string{"subdir1", "subdir2", "subdir3"}
	for i, dir := range subDirs {
		tDir := filepath.Join(tempDir, dir)
		tSubDir := filepath.Join(tDir, subSubDirs[i])
		expectedDirs = append(expectedDirs, tDir)
		expectedDirs = append(expectedDirs, tSubDir)
		err := os.Mkdir(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.Mkdir(tSubDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}
	}

	dirs, err := getAllDirs(tempDir)
	if err != nil {
		t.Fatalf("Failed to get directories: %v", err)
	}

	// Verify that the expected directories are present
	if len(dirs) != len(expectedDirs) {
		t.Errorf("Expected %d directories, but got %d", len(expectedDirs), len(dirs))
	}

	for _, dir := range expectedDirs {
		found := false
		for _, d := range dirs {
			if dir == d {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected directory %s not found", dir)
		}
	}
}

func Test_getAllFiles(t *testing.T) {
	tempDir := t.TempDir()

	expectedFiles := []string{"file1", "file2", "file3"}

	// Create some files inside the temporary directory
	for _, file := range expectedFiles {
		f, err := os.Create(filepath.Join(tempDir, file))
		f.Close()
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	files, err := getAllFiles(tempDir)
	if err != nil {
		t.Fatalf("Failed to get files: %v", err)
	}

	// Verify that the expected files are present
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(files))
	}

	for _, file := range expectedFiles {
		found := false
		for _, f := range files {
			if file == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found", file)
		}
	}
}
