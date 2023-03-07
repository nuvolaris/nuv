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
	"log"
	"os"
	"path/filepath"

	"github.com/nuvolaris/nuv/tools"
	"github.com/nuvolaris/task/cmd/taskmain/v3"
)

func main() {
	// initialize tools (used by the shell to find myself)
	tools.NuvCmd, _ = filepath.Abs(os.Args[0])

	// first argument with prefix "-" is an embedded tool
	// using "-" or "--" or "-task" invokes embedded task
	args := os.Args
	if len(args) > 1 && len(args[1]) > 0 && args[1][0] == '-' {
		cmd := args[1][1:]
		if cmd == "" || cmd == "-" || cmd == "task" {
			taskmain.Task(args[2:])
			os.Exit(0)
		}
		if cmd == "help" {
			tools.Help()
			os.Exit(0)
		}
		// check if it is an embedded to and invoke it
		if tools.IsTool(cmd) {
			code, err := tools.RunTool(cmd, args[2:])
			if err != nil {
				log.Print(err.Error())
			}
			os.Exit(code)
		}
		// no embeded tool found
		log.Printf("unknown tool -%s", cmd)
		os.Exit(0)
	}

	// execute nuv
	Nuv(getNuvRoot(), args[1:])
}
