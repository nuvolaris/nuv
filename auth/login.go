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

package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nuvolaris/nuv/config"
	"github.com/zalando/go-keyring"
)

type LoginResult struct {
	Login   string
	Auth    string
	ApiHost string
}

const usage = `Usage:
nuv login <apihost> [<user>]

Login to a Nuvolaris instance. If no user is specified, the default user "nuvolaris" is used.

Options:
  -h, --help   Show usage`

const whiskLoginPath = "/api/v1/web/whisk-system/nuv/login"
const defaultUser = "nuvolaris"
const nuvSecretServiceName = "nuvolaris"

func LoginCmd() (*LoginResult, error) {
	flag := flag.NewFlagSet("login", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Println(usage)
	}

	var helpFlag bool
	flag.BoolVar(&helpFlag, "h", false, "Show usage")
	flag.BoolVar(&helpFlag, "help", false, "Show usage")
	err := flag.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if helpFlag {
		flag.Usage()
		return nil, nil
	}

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		return nil, errors.New("missing apihost")
	}

	password := os.Getenv("NUV_PASSWORD")
	if password == "" {
		fmt.Print("Enter Password: ")
		pwd, err := AskPassword()
		if err != nil {
			return nil, err
		}
		password = pwd
		fmt.Println()
	}

	apihost := args[0]
	url := apihost + whiskLoginPath
	user := os.Getenv("NUV_LOGIN")
	if user == "" {
		user = defaultUser
	}
	if len(args) > 1 {
		user = args[1]
	}
	log.Println("Logging in as", user, "to", apihost)

	creds, err := doLogin(url, user, password)
	if err != nil {
		return nil, err
	}

	if _, ok := creds["AUTH"]; !ok {
		return nil, errors.New("missing AUTH token from login response")
	}

	nuvHome, err := homedir.Expand("~/.nuv")
	if err != nil {
		return nil, err
	}

	configMap, err := config.NewConfigMapBuilder().
		WithConfigJson(filepath.Join(nuvHome, "config.json")).
		Build()

	if err != nil {
		return nil, err
	}

	for k, v := range creds {
		if err := configMap.Insert(k, v); err != nil {
			return nil, err
		}
	}

	err = configMap.SaveConfig()
	if err != nil {
		return nil, err
	}

	// if err := storeCredentials(creds); err != nil {
	// 	return nil, err
	// }

	// auth, err := keyring.Get(nuvSecretServiceName, "AUTH")
	// if err != nil {
	// 	return nil, err
	// }

	return &LoginResult{
		Login:   user,
		Auth:    creds["AUTH"],
		ApiHost: apihost,
	}, nil
}

func doLogin(url, user, password string) (map[string]string, error) {
	data := map[string]string{
		"login":    user,
		"password": password,
	}
	loginJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("login failed with status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(body))
	}

	var creds map[string]string
	err = json.NewDecoder(resp.Body).Decode(&creds)
	if err != nil {
		return nil, errors.New("failed to decode response from login request")
	}

	return creds, nil
}

func storeCredentials(creds map[string]string) error {
	for k, v := range creds {
		err := keyring.Set(nuvSecretServiceName, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
