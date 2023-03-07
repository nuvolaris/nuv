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
)

func help() {
	if exists(".", NUVOPTS) {
		fmt.Println(readfile(NUVOPTS))
	} else {
		//fmt.Println("-t", "Nuvfile", "-l")
		Task("-t", NUVFILE, "-l")
	}
}

func parseArgs(args []string) []string {
	return args
}

// Nuv parse args moving
// into the folder corresponding to args
// then parse them with docopts and invokes task
func Nuv(base string, args []string) error {
	// go down using args as subcommands
	err := os.Chdir(base)
	if err != nil {
		return err
	}
	rest := args
	for _, dir := range args {
		if exists(dir, NUVFILE) {
			os.Chdir(dir)
			rest = rest[1:]
		} else {
			break
		}
	}

	if len(rest) == 0 || rest[0] == "help" {
		help()
		return nil
	}

	// parsed args
	if exists(".", NUVOPTS) {
		parsedArgs := parseArgs(rest)
		Task(parsedArgs...)
		return nil
	}
	// unparsed args
	Task(rest...)
	return nil
}
