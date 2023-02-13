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
package tools

import (
	"fmt"
	"os"

	"github.com/nojima/httpie-go"
)

func IsTool(name string) bool {
	if IsUtil(name) {
		return true
	}
	switch name {
	case "wsk":
		return true
	case "ht":
		return true
	}
	return false
}

func RunTool(name string, args []string) (int, error) {
	if IsUtil(name) {
		return RunUtil(name, args)
	}

	switch name {
	case "wsk":
		fmt.Println("=== wsk ===")
		cmd := append([]string{"wsk"}, args...)
		if err := Wsk(cmd); err != nil {
			return 1, err
		}
		return 0, nil
	case "ht":
		fmt.Println("=== ht ===")
		os.Args = append([]string{"ht"}, args...)
		if err := httpie.Main(); err != nil {
			return 1, err
		}
	}
	return 0, nil
}

func Help() {
	fmt.Println("Available tools:")
	tools := append(Utils, "wsk", "ht", "task")
	for _, x := range tools {
		fmt.Printf("-%s\n", x)
	}
}
