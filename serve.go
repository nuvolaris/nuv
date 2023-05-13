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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/browser"
)

var WebDir = "web"

const usage = `Usage:
nuv -serve [options]
-h, --help Print help message
--no-browser Do not open browser
`

// struct modeling this: { "stdout": ["stdout", "of", "nuv"], "stderr":  ["stderr", "of", "nuv"], "status": 0 }

type NuvOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Status int    `json:"status"`
}

func Serve(olarisDir string, args []string) error {
	// Define command line flags
	flag.Usage = func() { fmt.Print(usage) }

	var helpFlag bool
	var noBrowserFlag bool

	flag.BoolVar(&helpFlag, "h", false, "Print help message")
	flag.BoolVar(&helpFlag, "help", false, "Print help message")
	flag.BoolVar(&noBrowserFlag, "no-browser", false, "Do not open browser")

	// Parse command line flags
	os.Args = args
	flag.Parse()

	// Print help message if -h flag is provided
	if helpFlag {
		flag.Usage()
		return nil
	}

	// run nuv server and open browser
	port := getNuvPort()
	webDirPath := joinpath(olarisDir, WebDir)
	log.Println("Found web directory at: " + webDirPath)

	if !noBrowserFlag {
		if err := browser.OpenURL("http://localhost:" + port); err != nil {
			return err
		}
	}

	fileServer := webFileServerHandler(webDirPath)
	addr := fmt.Sprintf(":%s", port)

	http.Handle("/", fileServer)
	http.HandleFunc("/api/nuv", nuvTaskServer)

	if checkPortAvailable(port) {
		log.Println("Nuvolaris server started at http://localhost:" + port)
		return http.ListenAndServe(addr, nil)
	} else {
		log.Println("Nuvolaris server failed to start. Port is already in use.")
		return nil
	}
}

// Handler to serve the olaris/web directory
func webFileServerHandler(webDir string) http.Handler {
	return http.FileServer(http.Dir(webDir))
}

// The handler for the /api/ endpoint that runs nuv tasks
func nuvTaskServer(w http.ResponseWriter, r *http.Request) {
	trace("Request received:", r.URL.Path, r.URL.RawQuery)

	query := r.URL.RawQuery
	reqTasks := strings.Split(query, "+")

	nuvOut := execCommandTasks(reqTasks)

	// write output to response
	outjson, _ := json.Marshal(nuvOut)
	if _, err := w.Write(outjson); err != nil {
		debug("Failed to write response", err)
	}
}

func execCommandTasks(tasks []string) NuvOutput {
	if taskDryRun {
		return NuvOutput{
			Stdout: "Dry run: task " + strings.Join(tasks, " "),
			Stderr: "",
			Status: 0,
		}
	}
	trace("Running tasks from api:", tasks)

	cmd := exec.Command("nuv", tasks...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// We are using cmd.ProcessState.ExitCode() to get the exit code
	err := cmd.Run()
	if err != nil {
		warn("Failed to run tasks", tasks, stderr.String(), err)
	}
	return NuvOutput{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
		Status: cmd.ProcessState.ExitCode(),
	}

}

func checkPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	//nolint:errcheck
	ln.Close()
	return true
}
