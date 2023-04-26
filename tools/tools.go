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

	gojq "github.com/itchyny/gojq/cli"
	"github.com/nojima/httpie-go"
	envsubst "github.com/nuvolaris/envsubst/cmd/envsubstmain"
	"github.com/nuvolaris/goawk"
)

var tools = []string{
	"awk", "jq", "js", "envsubst", "wsk", "ht", "mkdir", "filetype", "random",
}

func availableCmds() []string {
	cmds := append(Utils, tools...)
	extra_cmds := []string{"config", "update", "serve", "help", "info", "version", "retry", "task"}
	cmds = append(cmds, extra_cmds...)
	return cmds
}

func IsTool(name string) bool {
	if IsUtil(name) {
		return true
	}
	for _, s := range tools {
		if s == name {
			return true
		}
	}
	return false
}

var NuvCmd = ""

func GetNuvCmd() string {
	if NuvCmd != "" {
		return NuvCmd
	}
	nuv := os.Getenv("NUV_CMD")
	if nuv != "" {
		return nuv
	}
	return ""
}

func RunTool(name string, args []string) (int, error) {
	if IsUtil(name) {
		return RunUtil(name, args)
	}
	switch name {
	case "wsk":
		//fmt.Println("=== wsk ===")
		cmd := append([]string{"wsk"}, args...)
		if err := Wsk(cmd); err != nil {
			return 1, err
		}
	case "ht":
		//fmt.Println("=== ht ===")
		os.Args = append([]string{"ht"}, args...)
		if err := httpie.Main(); err != nil {
			return 1, err
		}
	case "awk":
		// fmt.Println("=== awk ===")
		os.Args = append([]string{"goawk"}, args...)
		if err := goawk.AwkMain(); err != nil {
			return 1, err
		}
	case "jq":
		os.Args = append([]string{"gojq"}, args...)
		return gojq.Run(), nil
	case "js":
		os.Args = append([]string{"goja"}, args...)
		if err := jsToolMain(); err != nil {
			return 1, err
		}
	case "envsubst":
		os.Args = append([]string{"envsubst"}, args...)
		if err := envsubst.EnvsubstMain(); err != nil {
			return 1, err
		}
	case "mkdir":
		os.Args = append([]string{"mkdir"}, args...)
		if err := Mkdirs(); err != nil {
			return 1, err
		}
	case "filetype":
		os.Args = append([]string{"mkdir"}, args...)
		if err := Filetype(); err != nil {
			return 1, err
		}
	case "random":
		os.Args = append([]string{"random"}, args...)
		if err := RandTool(); err != nil {
			return 1, err
		}
	}

	return 0, nil
}

func Help() {
	fmt.Println("Available tools:")
	for _, x := range availableCmds() {
		fmt.Printf("-%s\n", x)
	}
}
