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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
)

var WebDir = "web"

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
	go RunNuvServer(webDirPath, port)

	if !noBrowserFlag {
		if err := browser.OpenURL("http://localhost:" + port); err != nil {
			return err
		}
	}

	select {}
}

func RunNuvServer(webDirPath string, port string) error {
	handler := nuvServerHandler(webDirPath)
	addr := fmt.Sprintf(":%s", port)

	log.Println("Nuvolaris server started at http://localhost:" + port)
	return http.ListenAndServe(addr, handler)
}

func nuvServerHandler(webDir string) http.Handler {
	return http.FileServer(http.Dir(webDir))
}

const usage = `Usage:
nuv -serve [options]
-h, --help Print help message
--no-browser Do not open browser
`
