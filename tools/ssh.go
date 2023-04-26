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
	"log"
	osuser "os/user"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

var (
	help bool

	user       string
	addr       string
	port       uint
	key        string
	cmd        string
	pass       bool
	passphrase bool
	timeout    time.Duration
	agent      bool

	err error
)

func SshTool() error {
	flag.Usage = printSSHUsage

	flag.BoolVar(&help, "help", false, "print this help message.")
	flag.BoolVar(&help, "h", false, "print this help message.")

	flag.StringVar(&user, "user", "root", "ssh user.")
	flag.StringVar(&addr, "ip", "127.0.0.1", "machine ip address.")
	flag.UintVar(&port, "port", 22, "ssh port number.")
	flag.StringVar(&key, "key", "", "private key path.")
	flag.StringVar(&cmd, "cmd", "", "command to run.")
	flag.BoolVar(&pass, "pass", false, "ask for ssh password instead of private key.")
	flag.BoolVar(&agent, "agent", false, "use ssh agent for authentication (unix systems only).")
	flag.BoolVar(&passphrase, "passphrase", false, "ask for private key passphrase.")
	flag.DurationVar(&timeout, "timeout", 0, "interrupt a command with SIGINT after a given timeout (0 means no timeout)")

	// Parse command line flags
	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	// if key is not provided, use default ~/.ssh/id_rsa
	if !isFlagPassed("key") {
		k, err := defaultSshKeyPath()
		if err != nil {
			return err
		}
		key = k
	}

	// if user is not provided, use current user
	if !isFlagPassed("user") {
		user = defaultSshUser()
	}

	return nil
}

func printSSHUsage() {
	fmt.Print(`Usage:
nuv -ssh [options]

-h, --help 		   print this help message.
-u, --user 		   the ssh user (default: current user)
`)
}

func defaultSshUser() string {
	usr, err := osuser.Current()
	if err != nil {
		log.Println("couldn't determine current user. Defaulting to 'root'")
		return "root"
	}
	return usr.Username
}

func defaultSshKeyPath() (string, error) {
	homessh, err := homedir.Expand("~/.ssh")
	if err != nil {
		return "", err
	}
	return filepath.Join(homessh, "id_rsa"), nil
}
