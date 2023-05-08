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
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	osuser "os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/nuvolaris/goph"
	"github.com/nuvolaris/nuv/auth"
	"golang.org/x/crypto/ssh"
)

func SshTool() error {
	var (
		help bool

		authMethod     goph.Auth
		user           string
		addr           string
		port           uint
		key            string
		cmd            string
		passwordFlag   bool
		passphraseFlag bool
		timeout        time.Duration
	)

	flag.Usage = printSSHUsage

	flag.BoolVar(&help, "help", false, "print this help message.")
	flag.BoolVar(&help, "h", false, "print this help message.")

	flag.StringVar(&user, "user", "root", "ssh user.")
	flag.StringVar(&user, "u", "root", "ssh user.")
	flag.UintVar(&port, "port", 22, "ssh port number.")
	flag.UintVar(&port, "p", 22, "ssh port number.")
	flag.StringVar(&key, "key", "", "private key path.")
	flag.StringVar(&key, "k", "", "private key path.")
	flag.BoolVar(&passwordFlag, "password", false, "ask for ssh password instead of private key.")
	flag.BoolVar(&passphraseFlag, "passphrase", false, "ask for ssh key passphrase.")
	flag.DurationVar(&timeout, "timeout", 0, "ssh connection timeout.")

	// Parse command line flags
	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	// if key is not provided, use default ~/.ssh/id_rsa
	if !isFlagPassed("key") && !isFlagPassed("k") {
		k, err := defaultSshKeyPath()
		if err != nil {
			return err
		}
		key = k
	}

	// if user is not provided, use current user
	if !isFlagPassed("user") && !isFlagPassed("u") {
		user = defaultSshUser()
	}

	var err error
	authMethod, err = getAuthMethod(key, passwordFlag, passphraseFlag)
	if err != nil {
		return err
	}

	// retrieve host from command line arguments
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return nil
	} else if len(args) == 1 {
		flag.Usage()
		return errors.New("missing command argument")
	}

	addr = args[0]
	cmd = strings.Join(args[1:], " ")

	client, err := goph.NewConn(&goph.Config{
		User:     user,
		Addr:     addr,
		Port:     port,
		Auth:     authMethod,
		Callback: VerifyHost,
	})

	if err != nil {
		return err
	}

	defer client.Close()

	// If the cmd flag passed
	if cmd != "" {
		ctx := context.Background()
		// create a context with timeout, if supplied in the argumetns
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		out, err := client.RunContext(ctx, cmd)

		fmt.Println(string(out), err)
		return err
	}

	return nil
}

func getAuthMethod(key string, askPassword, askPassphrase bool) (goph.Auth, error) {
	if goph.HasAgent() {
		return goph.UseAgent()
	} else if askPassword {
		askedPassword, err := askPass("Enter SSH Password: ")
		if err != nil {
			return nil, err
		}
		return goph.Password(askedPassword), nil
	} else if askPassphrase {
		askedPassphase, err := askPass("Enter Private Key Passphrase: ")
		if err != nil {
			return nil, err
		}
		return goph.Key(key, askedPassphase)
	}
	return nil, errors.New("no authentication method provided")
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

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	// if err != nil either key not in known hosts file OR host in known hosts file but key changed
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Key mismatch in this case (possible man in the middle attack)
	if hostFound && err != nil {
		return err
	}

	if hostFound && err == nil {
		return nil
	}

	// Ask user to check if he trust the host public key.
	if !askIsHostTrusted(host, key) {
		// Make sure to return error on non trusted keys.
		return errors.New("aborted")
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

func askIsHostTrusted(host string, key ssh.PublicKey) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key))
	fmt.Print("Would you like to add it? type yes or no: ")

	a, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(strings.TrimSpace(a)) == "yes"
}

func askPass(msg string) (string, error) {
	fmt.Print(msg)
	pass, err := auth.AskPassword()
	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(pass)), nil
}

func printSSHUsage() {
	fmt.Print(`Usage:
nuv -ssh [options] <address> <cmd> 

Connect to a remote host via ssh. If no command is provided, an interactive shell is opened.

Options:
-h, --help 		print this help message.

-u, --user STRING 		the ssh user (default: current user)
-k, --key STRING 		the private key path (default: ~/.ssh/id_rsa)
-p, --port INT 			the ssh port number (default: 22)
--password			ask for ssh password instead of private key
--passphrase			ask for ssh key passphrase (default: false)
--timeout INT			ssh connection timeout (default: 0, meaning no timeout)
`)
}
