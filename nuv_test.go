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
	"path/filepath"
)

func ExampleNuv() {
	// test
	olaris, _ := filepath.Abs("olaris")
	pr(1, Nuv(olaris, split("")))
	pr(2, Nuv(olaris, split("top")))
	pr(3, Nuv(olaris, split("top arg")))
	pr(4, Nuv(olaris, split("sub")))
	pr(5, Nuv(olaris, split("sub opts")))
	pr(6, Nuv(olaris, split("sub opts args 1")))
	// Output:
	// (olaris) task [-t nuvfile.yml -l]
	// 1 <nil>
	// (olaris) task [-t nuvfile.yml top --]
	// 2 <nil>
	// (olaris) task [-t nuvfile.yml top -- arg]
	// 3 <nil>
	// (sub) task [-t nuvfile.yml -l]
	// 4 <nil>
	// Usage:
	//   opts args <name>... [-c]
	//   opts arg1 <name> arg2 <x> <y> [--fl=<flag>]
	//   opts arg3 (opt1|opt2) <x> <y> [--fa|--fb]
	//   opts -h | --help | --version
	//
	// 5 <nil>
	// (opts) task [__fa=false __fb=false __fl= __help=false __version=false _c=false _h=false _name_=('1') _x_= _y_= arg1=false arg2=false arg3=false args=true opt1=false opt2=false]
	// 6 <nil>
}

func ExampleParseArgs() {
	usage := readfile("olaris/sub/opts/nuvopts.txt")

	pr(1, parseArgs(usage, split("args mike miri max")))
	pr(2, parseArgs(usage, split("args mike -c")))
	pr(3, parseArgs(usage, split("arg1 max arg2 1 2 --fl=3")))
	pr(4, parseArgs(usage, split("arg3 opt2 4 5 --fb")))

	// Output:
	// 1 [__fa=false __fb=false __fl= __help=false __version=false _c=false _h=false _name_=('mike' 'miri' 'max') _x_= _y_= arg1=false arg2=false arg3=false args=true opt1=false opt2=false]
	// 2 [__fa=false __fb=false __fl= __help=false __version=false _c=true _h=false _name_=('mike') _x_= _y_= arg1=false arg2=false arg3=false args=true opt1=false opt2=false]
	// 3 [__fa=false __fb=false __fl=3 __help=false __version=false _c=false _h=false _name_=('max') _x_=1 _y_=2 arg1=true arg2=true arg3=false args=false opt1=false opt2=false]
	// 4 [__fa=false __fb=true __fl= __help=false __version=false _c=false _h=false _name_=() _x_=4 _y_=5 arg1=false arg2=false arg3=true args=false opt1=false opt2=true]
}
