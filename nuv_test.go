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
	"os"
	"path/filepath"
)

func ExampleNuv() {
	// test
	olaris, _ := filepath.Abs("olaris")
	pr(1, Nuv(olaris, split("")))
	pr(2, Nuv(olaris, split("sub")))
	pr(3, Nuv(olaris, split("sub multi")))
	pr(4, Nuv(olaris, split("sub multi ship")))

	// Output:
	// -
}

func ExampleNuvCmd() {
	os.Chdir(homeDir)
	pr(1, Nuv("tests", as("hello")))
	// Output:
	// -
}

func ExampleTask() {
	Task("task", "-l")
	// Output:
	// -
}
