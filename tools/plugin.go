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
	"flag"
	"fmt"
	"os"
)

func pluginTool() error {

	flag := flag.NewFlagSet("plugin", flag.ExitOnError)
	flag.Usage = printPluginUsage

	err := flag.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return nil
	}

	fmt.Println("plugin tool not implemented yet")

	return nil
}

func printPluginUsage() {
	fmt.Println(`Usage: nuv -plugin repo

Install/update plugins from a repository.
The name of the repository must start with 'olaris-'.`)
}
