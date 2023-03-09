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
	"sort"
	"strings"

	docopt "github.com/docopt/docopt-go"
)

func help() {
	if exists(".", NUVOPTS) {
		fmt.Println(readfile(NUVOPTS))
	} else {
		//fmt.Println("-t", "Nuvfile", "-l")
		Task("-t", NUVFILE, "-l")
	}
}

// parseArgs parse the arguments acording the docopt
// it returns a sequence suitable to be feed as arguments for task.
// note that it will change hyphens for flags ('-c', '--count') to '_' ('_c' '__count')
// and '<' and '>' for parameters '_' (<hosts> => _hosts_)
// boolean are "true" or "false" and arrays in the form ('first' 'second')
// suitable to be used as arrays
// Examples:
// if "Usage: nettool ping [--count=<max>] <hosts>..."
// with "ping --count=3 google apple" returns
// ping=true _count=3 _hosts_=('google' 'apple')
func parseArgs(usage string, args []string) []string {
	res := []string{}
	// parse args
	parser := docopt.Parser{}
	opts, err := parser.ParseArgs(usage, args, NuvVersion)
	if err != nil {
		return res
	}
	for k, v := range opts {
		kk := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(k, "-", "_"), "<", "_"), ">", "_")
		vv := ""
		//fmt.Println(v, reflect.TypeOf(v))
		switch o := v.(type) {
		case bool:
			vv = "false"
			if o {
				vv = "true"
			}
		case string:
			vv = o
		case []string:
			a := []string{}
			for _, i := range o {
				a = append(a, fmt.Sprintf("'%v'", i))
			}
			vv = "(" + strings.Join(a, " ") + ")"
		case nil:
			vv = ""
		}
		res = append(res, fmt.Sprintf("%s=%s", kk, vv))
	}
	sort.Strings(res)
	return res
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
			//fmt.Println(dir, rest)
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
		fmt.Println("PREPARSE:", rest)
		parsedArgs := parseArgs(readfile(NUVOPTS), rest)
		prefix := []string{"-t", NUVFILE}
		if len(rest) > 0 && rest[0][0] != '-' {
			prefix = append(prefix, rest[0])
		}
		parsedArgs = append(prefix, parsedArgs...)
		fmt.Println("POSTPARSE:", parsedArgs)
		pwd, _ := os.Getwd()
		fmt.Println("PWD", pwd)
		Task(parsedArgs...)
		return nil
	}
	// unparsed args
	taskArgs := append([]string{"-t", NUVFILE, rest[0], "--"}, rest[1:]...)
	Task(taskArgs...)
	return nil
}
