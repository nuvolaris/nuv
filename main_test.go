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
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

/// test utils

// pr print args
func pr(args ...any) {
	fmt.Println(args...)
}

// as creates a string array
// func as(s ...string) []string {
// 	return s
// }

var homeDir = ""
var workDir = ""

// normalize path replacing the variable path with /work
func npath(dir string) string {
	return strings.Replace(dir, workDir, "/work", -1)
}

// normalize path replaching the variable home part with /home
func nhpath(dir string) string {
	//fmt.Println("dir", dir, "home", homeDir)
	return strings.Replace(dir, homeDir, "/home", -1)
}

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	workDir, _ = filepath.Abs(wd)
	homeDir, _ = homedir.Dir()
	taskDryRun = true
	//debugging = true
	//tracing = true
	os.Exit(m.Run())
}

func TestSetupNuvRootPlugin(t *testing.T) {
	// Test case 1: NUV_ROOT_PLUGIN is not set
	os.Unsetenv("NUV_ROOT_PLUGIN")
	os.Setenv("NUV_PWD", "/path/to/nuv")
	setNuvRootPluginEnv()
	if os.Getenv("NUV_ROOT_PLUGIN") != "/path/to/nuv" {
		t.Errorf("NUV_ROOT_PLUGIN not set correctly, expected /path/to/nuv but got %s", os.Getenv("NUV_ROOT_PLUGIN"))
	}

	// Test case 2: NUV_ROOT_PLUGIN is already set
	os.Setenv("NUV_ROOT_PLUGIN", "/path/to/nuv/root")
	setNuvRootPluginEnv()
	if os.Getenv("NUV_ROOT_PLUGIN") != "/path/to/nuv/root" {
		t.Errorf("NUV_ROOT_PLUGIN not set correctly, expected /path/to/nuv/root but got %s", os.Getenv("NUV_ROOT_PLUGIN"))
	}
}
