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
	"testing"
)

func ExampleNuvArg() {
	// test
	os.Chdir(workDir)
	olaris, _ := filepath.Abs(joinpath("tests", "olaris"))
	err := Nuv(olaris, split("top"))
	pr(2, err)
	err = Nuv(olaris, split("top arg"))
	pr(3, err)
	err = Nuv(olaris, split("top arg VAR=1"))
	pr(4, err)
	err = Nuv(olaris, split("top VAR=1 arg"))
	pr(5, err)
	// Output:
	// (olaris) task [-t nuvfile.yml top --]
	// 2 <nil>
	// (olaris) task [-t nuvfile.yml top -- arg]
	// 3 <nil>
	// (olaris) task [-t nuvfile.yml top VAR=1 -- arg]
	// 4 <nil>
	// (olaris) task [-t nuvfile.yml top VAR=1 -- arg]
	//5 <nil>
}

func ExampleNuv() {
	// test
	os.Chdir(workDir)
	olaris, _ := filepath.Abs(joinpath("tests", "olaris"))
	err := Nuv(olaris, split(""))
	pr(1, err)
	err = Nuv(olaris, split("sub"))
	pr(4, err)
	err = Nuv(olaris, split("sub opts"))
	pr(5, err)
	err = Nuv(olaris, split("sub opts args 1"))
	pr(6, err)
	// Output:
	// (olaris) task [-t nuvfile.yml -l]
	// 1 <nil>
	// (sub) task [-t nuvfile.yml -l]
	// 4 <nil>
	// Usage:
	//   opts hello
	//   opts args <name>... [-c]
	//   opts arg1 <name> arg2 <x> <y> [--fl=<flag>]
	//   opts arg3 (opt1|opt2) <x> <y> [--fa|--fb]
	//   opts -h | --help | --version
	//
	// 5 <nil>
	// (opts) task [-t nuvfile.yml args __fa=false __fb=false __fl= __help=false __version=false _c=false _h=false _name_=('1') _x_= _y_= arg1=false arg2=false arg3=false args=true hello=false opt1=false opt2=false]
	// 6 <nil>
}

func ExampleParseArgs() {
	os.Chdir(workDir)
	usage := readfile("tests/olaris/sub/opts/nuvopts.txt")
	args := parseArgs(usage, split("args mike miri max"))
	pr(1, args)
	args = parseArgs(usage, split("args mike -c"))
	pr(2, args)
	args = parseArgs(usage, split("arg1 max arg2 1 2 --fl=3"))
	pr(3, args)
	args = parseArgs(usage, split("arg3 opt2 4 5 --fb"))
	pr(4, args)

	// Output:
	// 1 [__fa=false __fb=false __fl= __help=false __version=false _c=false _h=false _name_=('mike' 'miri' 'max') _x_= _y_= arg1=false arg2=false arg3=false args=true hello=false opt1=false opt2=false]
	// 2 [__fa=false __fb=false __fl= __help=false __version=false _c=true _h=false _name_=('mike') _x_= _y_= arg1=false arg2=false arg3=false args=true hello=false opt1=false opt2=false]
	// 3 [__fa=false __fb=false __fl=3 __help=false __version=false _c=false _h=false _name_=('max') _x_=1 _y_=2 arg1=true arg2=true arg3=false args=false hello=false opt1=false opt2=false]
	// 4 [__fa=false __fb=true __fl= __help=false __version=false _c=false _h=false _name_=() _x_=4 _y_=5 arg1=false arg2=false arg3=true args=false hello=false opt1=false opt2=true]
}

func Test_getTaskNamesList(t *testing.T) {
	t.Run("empty nuvfile should return empty array", func(t *testing.T) {
		tmpDir := createTmpNuvfile(t, "")

		tasks := getTaskNamesList(tmpDir)
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("should return array of task names if tasks in nuvfile", func(t *testing.T) {
		tmpDir := createTmpNuvfile(t, "tasks:\n  task1: a\n  task2: b\n")
		defer os.RemoveAll(tmpDir)

		tasks := getTaskNamesList(tmpDir)
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}

		if tasks[0] != "task1" || tasks[1] != "task2" {
			t.Fatalf("expected task1 and task2, got %s and %s", tasks[0], tasks[1])
		}
	})

}

func createTmpNuvfile(t *testing.T, content string) string {
	t.Helper()
	// create temp folder with nuvfile.yml
	tmpDir, err := os.MkdirTemp("", "nuv-test")
	if err != nil {
		t.Fatal(err)
	}

	// create nuvfile.yml
	nuvfile := filepath.Join(tmpDir, "nuvfile.yml")
	err = os.WriteFile(nuvfile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return tmpDir
}
