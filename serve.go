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
	"net"
	"net/http"
	"os"

	"github.com/pkg/browser"
)

func Serve(olarisDir string, args []string) error {
	flag := flag.NewFlagSet("serve", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Println(`Serve a local directory on http://localhost:9768. You can change port with the NUV_PORT environment variable.

Usage:
  nuv -serve [options] <dir>

Options:
  -h, --help Print help message
  --no-open Do not open browser automatically
  --proxy <proxy> Use proxy server
		`)
	}
	// Define command line flags

	var helpFlag bool
	var noBrowserFlag bool
	var proxyFlag string

	flag.BoolVar(&helpFlag, "h", false, "Print help message")
	flag.BoolVar(&helpFlag, "help", false, "Print help message")
	flag.BoolVar(&noBrowserFlag, "no-open", false, "Do not open browser")
	flag.StringVar(&proxyFlag, "proxy", "", "Use proxy server")

	// Parse command line flags
	os.Args = args
	err := flag.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	// Print help message if requested
	if flag.NArg() != 1 || helpFlag {
		flag.Usage()
		return nil
	}

	webDir := flag.Arg(0)

	// run nuv server and open browser
	port := getNuvPort()
	webDirPath := joinpath(os.Getenv("NUV_PWD"), webDir)
	log.Println("Serving directory: " + webDirPath)

	if !noBrowserFlag {
		if err := browser.OpenURL("http://localhost:" + port); err != nil {
			return err
		}
	}

	fileServer := webFileServerHandler(webDirPath)
	addr := fmt.Sprintf(":%s", port)

	http.Handle("/", fileServer)

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

func checkPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	//nolint:errcheck
	ln.Close()
	return true
}
