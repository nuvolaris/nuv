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
//
package main

import (
	"errors"
	"log"
	"os"
)

// Nuv enter in the directly and parse args
// It recurses if the first arg is a subdirectory
func Nuv(dir string, args []string) error {
	log.Printf("Nuv: %s %v", dir, args)
	// check what is the the dir
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New(dir + " should be a directory!")
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		Task("-l")
		return nil
	}
	// now process args and recurse if a directory
	cmd := args[0]
	fi, err = os.Stat(cmd)
	if os.IsNotExist(err) {
		return processCmd(args)
	}
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return Nuv(cmd, args)
	}
	return errors.New(dir + " exists and is not a directory")
}

// processCmd assumes cmd is a string and it is not a directory
func processCmd(args []string) error {
	log.Printf("ProcessCmd: %v", args)
	taskDryRun = true
	Task(args...)
	taskDryRun = false
	return nil
}
