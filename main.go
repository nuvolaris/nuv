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

	httpie "github.com/nojima/httpie-go"
	"github.com/nuvolaris/task/cmd/taskmain/v3"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("nuv, the next generation.\n")
	} else if args[1][0] == '-' {
		switch args[1] {
		case "-task", "-t":
			fmt.Println("task")
			args := append([]string{"task"}, args[2:]...)
			taskmain.Task(args)
			return
		case "-wsk", "-w":
			fmt.Println("wsk")
			args := append([]string{"task"}, args[2:]...)
			Wsk(args)
			return
		case "-ht", "-h":
			fmt.Println("ht")
			os.Args = append([]string{"task"}, args[2:]...)
			if err := httpie.Main(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Println("unknown")
		}

	}
}
