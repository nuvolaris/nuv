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
	"regexp"

	"github.com/jaytaylor/go-find"
)

func FindTool() error {
	var (
		help      bool
		name      string
		wholeName string
		maxDepth  int
		minDepth  int
		empty     bool
		typeStr   string
		regex     string
	)

	flag.Usage = printFindUsage

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&name, "n", "", "Name pattern to match")
	flag.StringVar(&name, "name", "", "Name pattern to match")
	flag.StringVar(&wholeName, "w", "", "Name pattern to match")
	flag.StringVar(&wholeName, "whole-name", "", "Name pattern to match")
	flag.IntVar(&maxDepth, "max-depth", 0, "Maximum depth to search")
	flag.IntVar(&minDepth, "min-depth", 0, "Minimum depth to search")
	flag.BoolVar(&empty, "e", false, "Match empty files and directories")
	flag.BoolVar(&empty, "empty", false, "Match empty files and directories")
	flag.StringVar(&typeStr, "t", "", "File type to match")
	flag.StringVar(&typeStr, "type", "", "File type to match")
	flag.StringVar(&regex, "r", "", "Regex pattern to match")
	flag.StringVar(&regex, "regex", "", "Regex pattern to match")

	// Parse command-line flags
	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	// Get paths args
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	if name == "" && wholeName == "" && regex == "" {
		flag.Usage()
		return nil
	}

	finder := find.NewFind(args...)

	if name != "" {
		finder = finder.Name(name)
	}

	if wholeName != "" {
		finder = finder.WholeName(wholeName)
	}

	if maxDepth != 0 {
		finder = finder.MaxDepth(maxDepth)
	}

	if minDepth != 0 {
		finder = finder.MinDepth(minDepth)
	}

	if empty {
		finder = finder.Empty()
	}

	if typeStr != "" {
		finder = finder.Type(typeStr)
	}

	if regex != "" {
		regexp, err := regexp.Compile(regex)
		if err != nil {
			return err
		}
		finder = finder.Regex(regexp)
	}

	results, err := finder.Evaluate()
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func printFindUsage() {
	fmt.Print(`Usage: nuv -find [options] [paths...]

Options:
-h, --help			Show help
-n, --name			Name pattern to match
-w, --whole-name		Name pattern to match
--max-depth			Maximum depth to search
--min-depth			Minimum depth to search
-e, --empty			Match empty files and directories
-t, --type			File type to match
-r, --regex			Regex pattern to match
`)
}
