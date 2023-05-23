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
	"regexp"
	"strconv"
)

func validateTool() error {
	flag := flag.NewFlagSet("validate", flag.ContinueOnError)

	flag.Usage = printValidateUsage

	helpHalf := flag.Bool("h", false, "")
	envFlag := flag.Bool("e", false, "")
	mailFlag := flag.Bool("m", false, "")
	numberFlag := flag.Bool("n", false, "")
	regexFlag := flag.String("r", "", "")

	err := flag.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if *helpHalf {
		flag.Usage()
		return nil
	}

	if flag.NArg() != 1 {
		flag.Usage()
		return nil
	}

	arg := flag.Arg(0)
	value := arg

	if *envFlag {
		value = os.Getenv(arg)
		if value == "" {
			fmt.Printf("The variable '%s' is not set.\n", arg)
			return nil
		}
	}

	if *mailFlag {
		if isValidEmail(value) {
			fmt.Printf("'%s' %sis a valid email address.\n", value, envVarMsg(*envFlag, arg))
		} else {
			fmt.Printf("'%s' %sis NOT a valid email address.\n", value, envVarMsg(*envFlag, arg))
		}
		return nil
	}

	if *numberFlag {
		if isValidNumber(value) {
			fmt.Printf("'%s' %sis a valid number.\n", value, envVarMsg(*envFlag, arg))
		} else {
			fmt.Printf("'%s' %sis NOT a valid number.\n", value, envVarMsg(*envFlag, arg))
		}
		return nil
	}

	if *regexFlag != "" {
		valid, err := isValidByRegex(value, *regexFlag)
		if err != nil {
			return err
		}

		if valid {
			fmt.Printf("'%s' %smatches the regex.\n", value, envVarMsg(*envFlag, arg))
		} else {
			fmt.Printf("'%s' %sdoes NOT match the regex.\n", value, envVarMsg(*envFlag, arg))
		}
		return nil
	}

	return nil
}

func envVarMsg(envFlag bool, name string) string {
	if envFlag {
		return "from the variable '" + name + "' "
	}
	return ""
}

func isValidNumber(number string) bool {
	_, err := strconv.ParseFloat(number, 64)
	return err == nil
}

func isValidEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Create a regular expression object
	regExp := regexp.MustCompile(pattern)

	// Use the regular expression to match the email string
	return regExp.MatchString(email)
}

func isValidByRegex(value string, regex string) (bool, error) {
	// Create a regular expression object
	regExp, err := regexp.Compile(regex)
	if err != nil {
		return false, err
	}

	// Use the regular expression to match the email string
	return regExp.MatchString(value), nil
}

func printValidateUsage() {
	fmt.Println(`Usage:
nuv -validate [-e] [-m | -n | -r <regex>] <value>

Check if a value is valid according to the given constraints.
If -e is specified, the value is retrieved from the environment variable with the given name.

Options:
	-h		Print this help message.
	-e		The value is retrieved from the environment variable with the given name.
	-m		Check if the value is a valid email address.
	-n		Check if the value is a number.
	-r		Check if the value matches the given regular expression.`)
}
