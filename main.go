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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nuvolaris/nuv/tools"
	"github.com/nuvolaris/task/cmd/taskmain/v3"
	"github.com/pkg/browser"
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
	root, _ := getNuvRoot()
	fmt.Println("ROOT:", root)
}

func main() {
	if os.Getenv("NUV_NO_LOG_PREFIX") != "" {
		log.SetFlags(0)
	}
	if os.Getenv("NUV_VERSION") != "" {
		NuvVersion = os.Getenv("NUV_VERSION")
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
		if cmd == "version" {
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
		if cmd == "update" {
			// ok no up, nor down, let's download it
			err := pullTasks(true, true)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			fmt.Println("nuvfiles updated successfully")
			os.Exit(0)
		}
		if cmd == "retry" {
			if err := tools.ExpBackoffRetry(args[1:]); err != nil {
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

	dir, err := getNuvRoot()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// check if olaris was recently updated
	// we pass parent(dir) because we use the olaris parent folder
	checkUpdated(parent(dir), 24*time.Hour)

	if len(args) == 1 {
		// run nuv server and open browser
		port := getNuvPort()
		go RunNuvServer(dir, port)
		if err := browser.OpenURL("http://localhost:" + port); err != nil {
			log.Fatal(err)
		}
		select {}
	} else {
		if err := Nuv(dir, args[1:]); err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	}
}
