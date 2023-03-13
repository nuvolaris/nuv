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

	"github.com/mitchellh/go-homedir"
)

func Example_locate() {
	dir, err := locateNuvRoot(".")
	pr(1, err, npath(dir))
	dir, err = locateNuvRoot("olaris")
	pr(2, err, npath(dir))
	dir, err = locateNuvRoot(joinpath("olaris", "sub"))
	pr(3, err, npath(dir))
	// Output:
	// 1 <nil> /work/olaris
	// 2 <nil> /work/olaris
	// 3 <nil> /work/olaris
}

func Example_locate_git() {
	NuvBranch = "test"
	nuvdir, _ := homedir.Expand("~/.nuv")
	os.RemoveAll(nuvdir)
	os.Setenv("NUV_BIN", "")
	//_, err :=
	_, err := locateNuvRoot("tests")
	pr(1, err)
	os.Setenv("NUV_BIN", workDir)
	dir, err := locateNuvRoot("tests")
	pr(2, err, npath(dir))
	downloadTasksFromGitHub(true, true)
	dir, err = locateNuvRoot("tests")
	pr(3, err, nhpath(dir))
	downloadTasksFromGitHub(true, true)
	dir, err = locateNuvRoot("tests")
	pr(4, err, nhpath(dir))
	// Output:
	// 1 we cannot find nuvfiles, download them with nuv -update
	// 2 <nil> /work/olaris
	// Cloning tasks...
	// 3 <nil> /home/.nuv/olaris
	// Updating tasks...
	// 4 <nil> /home/.nuv/olaris
}
