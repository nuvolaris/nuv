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
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nuvolaris/task/v3"
)

func printInstalledPluginsMessage(localDir string) error {
	plgs, err := newPlugins(localDir)
	if err != nil {
		return err
	}
	plgs.print()
	return nil
}

func findTaskInPlugins(localDir string, plg string) (string, error) {
	trace("findTaskInPlugins", plg)

	plgs, err := newPlugins(localDir)
	if err != nil {
		return "", err
	}

	for _, folder := range plgs.local {
		_, err := validateTaskName(folder, plg)
		if err == nil {
			return folder, nil
		}
	}

	return "", &TaskNotFoundErr{input: plg}
}

type plugins struct {
	local []string
	nuv   []string
}

func newPlugins(localDir string) (*plugins, error) {
	localOlarisFolders := make([]string, 0)
	nuvOlarisFolders := make([]string, 0)

	// Search in directory (localDir/olaris-*)
	dir := filepath.Join(localDir, "olaris-*")
	olarisFolders, err := filepath.Glob(dir)
	if err != nil {
		return nil, err
	}

	localOlarisFolders = append(localOlarisFolders, olarisFolders...)

	// Search in ~/.nuv/olaris-*
	nuvHome, err := homedir.Expand("~/.nuv")
	if err != nil {
		return nil, err
	}

	olarisNuvFolders, err := filepath.Glob(filepath.Join(nuvHome, "olaris-*"))
	if err != nil {
		return nil, err
	}

	nuvOlarisFolders = append(nuvOlarisFolders, olarisNuvFolders...)

	return &plugins{
		local: localOlarisFolders,
		nuv:   nuvOlarisFolders,
	}, nil
}

func (p *plugins) print() {
	if len(p.local) == 0 && len(p.nuv) == 0 {
		fmt.Println("No plugins installed. Use 'nuv -plugin' to add new ones.")
		return
	}

	fmt.Println("Plugins:")
	if len(p.local) > 0 {
		for _, plg := range p.local {
			fmt.Printf("[LOCAL] %s:\n", filepath.Base(plg))
			printTaskHelp(plg)
		}
	}

	if len(p.nuv) > 0 {
		for _, plg := range p.nuv {
			fmt.Printf("[NUV]  %s:\n", filepath.Base(plg))
			printTaskHelp(plg)
		}
	}
}

func printTaskHelp(path string) {
	dir := path
	entrypoint := NUVFILE
	e := task.Executor{
		Dir:        dir,
		Entrypoint: entrypoint,
		Summary:    true,
		Color:      true,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	e.ListTaskNames(true)
}
