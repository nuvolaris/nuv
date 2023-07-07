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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/nuvolaris/nuv/auth"
	"github.com/nuvolaris/nuv/config"
	"github.com/nuvolaris/nuv/tools"
	"github.com/nuvolaris/task/cmd/taskmain/v3"
)

func setupCmd(me string) (string, error) {
	if os.Getenv("NUV_CMD") != "" {
		return os.Getenv("NUV_CMD"), nil
	}

	// look in path
	me, err := exec.LookPath(me)
	if err != nil {
		return "", err
	}
	trace("found", me)

	// resolve links
	fileInfo, err := os.Lstat(me)
	if err != nil {
		return "", err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		me, err = os.Readlink(me)
		if err != nil {
			return "", err
		}
		trace("resolving link to", me)
	}

	// get the absolute path
	me, err = filepath.Abs(me)
	if err != nil {
		return "", err
	}
	trace("ME:", me)
	//nolint:errcheck
	os.Setenv("NUV_CMD", me)
	return me, nil
}

func setupBinPath(cmd string) {
	// initialize tools (used by the shell to find myself)
	if os.Getenv("NUV_BIN") == "" {
		os.Setenv("NUV_BIN", filepath.Dir(cmd))
	}
	os.Setenv("PATH", fmt.Sprintf("%s%c%s", os.Getenv("NUV_BIN"), os.PathListSeparator, os.Getenv("PATH")))
	debugf("PATH=%s", os.Getenv("PATH"))

	//subpath := fmt.Sprintf("\"%s\"%c\"%s\"", os.Getenv("NUV_BIN"), os.PathListSeparator, joinpath(os.Getenv("NUV_BIN"), runtime.GOOS+"-"+runtime.GOARCH))
	//os.Setenv("PATH", fmt.Sprintf("%s%c%s", subpath, os.PathListSeparator, os.Getenv("PATH")))
}

func info() {
	fmt.Println("VERSION:", NuvVersion)
	fmt.Println("BRANCH:", NuvBranch)
	fmt.Println("CMD:", tools.GetNuvCmd())
	fmt.Println("BIN:", os.Getenv("NUV_BIN"))
	fmt.Println("REPO:", getNuvRepo())
	fmt.Println("TMP:", os.Getenv("NUV_TMP"))
	root, _ := getNuvRoot()
	fmt.Println("ROOT:", root)
	fmt.Println("NUV_PWD:", os.Getenv("NUV_PWD"))
}

func main() {
	// set version
	if os.Getenv("NUV_VERSION") != "" {
		NuvVersion = os.Getenv("NUV_VERSION")
	}

	// disable log
	if os.Getenv("NUV_NO_LOG_PREFIX") != "" {
		log.SetFlags(0)
	}

	if pwd, err := os.Getwd(); err != nil {
		warn("unable to set NUV_PWD to working directory", err)
	} else {
		os.Setenv("NUV_PWD", pwd)
	}

	var err error
	me := os.Args[0]
	if filepath.Base(me) == "nuv" || filepath.Base(me) == "nuv.exe" {
		tools.NuvCmd, err = setupCmd(me)
		if err != nil {
			warn("cannot setup cmd", err)
			os.Exit(1)
		}
		setupBinPath(tools.NuvCmd)
	}

	nuvHome, err := homedir.Expand("~/.nuv")
	if err != nil {
		warn("cannot expand home dir", err)
		os.Exit(1)
	}

	// first argument with prefix "-" is an embedded tool
	// using "-" or "--" or "-task" invokes embedded task
	args := os.Args
	if len(args) > 1 && len(args[1]) > 0 && args[1][0] == '-' {
		cmd := args[1][1:]
		if cmd == "" || cmd == "-" || cmd == "task" {
			params := []string{"task"}
			if len(args) > 2 {
				params = append(params, args[2:]...)
			}
			exitCode, err := taskmain.Task(params)
			if err != nil {
				log.Println(err)
			}
			os.Exit(exitCode)
		}
		if cmd == "version" || cmd == "v" {
			fmt.Println(NuvVersion)
			os.Exit(0)
		}
		if cmd == "info" {
			info()
			os.Exit(0)
		}
		if cmd == "help" {
			tools.Help()
			os.Exit(0)
		}
		if cmd == "serve" {
			if err := Serve(retrieveRootDir(), args[1:]); err != nil {
				log.Fatalf("error: %v", err)
			}
			os.Exit(0)
		}
		if cmd == "update" {
			// ok no up, nor down, let's download it
			err := pullTasks(true, true)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		if cmd == "retry" {
			if err := tools.ExpBackoffRetry(args[1:]); err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			os.Exit(0)
		}
		if cmd == "login" {
			os.Args = args[1:]
			loginResult, err := auth.LoginCmd()
			if err != nil {
				log.Fatalf("error: %s", err.Error())
			}

			if loginResult == nil {
				os.Exit(1)
			}

			fmt.Println("Successfully logged in as " + loginResult.Login + ".")
			if err := wskPropertySet(loginResult.ApiHost, loginResult.Auth); err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			fmt.Println("Nuvolaris host and auth set successfully. You are now ready to use nuv -wsk!")
			os.Exit(0)
		}
		if cmd == "config" {
			os.Args = args[1:]
			nuvRootPath := joinpath(retrieveRootDir(), NUVROOT)
			configPath := joinpath(nuvHome, CONFIGFILE)
			if err := config.ConfigTool(nuvRootPath, configPath); err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			os.Exit(0)
		}
		if cmd == "plugin" {
			os.Args = args[1:]
			if err := pluginTool(); err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			os.Exit(0)
		}
		// check if it is an embedded to and invoke it
		if tools.IsTool(cmd) {
			code, err := tools.RunTool(cmd, args[2:])
			if err != nil {
				log.Print(err.Error())
			}
			os.Exit(code)
		}
		// no embeded tool found
		warn("unknown tool", "-"+cmd)
		os.Exit(0)
	}

	nuvRootDir := retrieveRootDir()

	setupTmp()

	err = setAllConfigEnvVars(nuvRootDir, nuvHome)
	if err != nil {
		warn("cannot apply env vars from configs", err)
		os.Exit(1)
	}

	// check if olaris was recently updated
	// we pass parent(dir) because we use the olaris parent folder
	checkUpdated(parent(nuvRootDir), 24*time.Hour)

	if err := runNuv(nuvRootDir, args); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

func retrieveRootDir() string {
	dir, err := getNuvRoot()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
	return dir
}

func setAllConfigEnvVars(nuvRootDir string, configDir string) error {
	trace("setting all config env vars")
	configMap, err := config.NewConfigMapBuilder().
		WithNuvRoot(joinpath(nuvRootDir, NUVROOT)).
		WithConfigJson(joinpath(configDir, CONFIGFILE)).
		Build()

	if err != nil {
		return err
	}

	kv := configMap.Flatten()
	for k, v := range kv {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
		debug("env var set", k, v)
	}

	return nil
}

func wskPropertySet(apihost, auth string) error {
	args := []string{"property", "set", "--apihost", apihost, "--auth", auth}
	cmd := append([]string{"wsk"}, args...)
	if err := tools.Wsk(cmd); err != nil {
		return err
	}
	return nil
}

func runNuv(baseDir string, args []string) error {
	err := Nuv(baseDir, args[1:])
	if err == nil {
		return nil
	}

	var taskNotFoundErr *TaskNotFoundErr
	if errors.As(err, &taskNotFoundErr) {
		// Hook plugins here
		trace("task not found, looking for plugins")
		plgDir, err := findTaskInPlugins(parent(baseDir), args[1])
		if err != nil {
			return taskNotFoundErr
		}

		debug("Found plugin in", plgDir)
		os.Setenv("NUV_ROOT", plgDir)
		if err := Nuv(plgDir, args[1:]); err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		return nil
	}

	return err
}
