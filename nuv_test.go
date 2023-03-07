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
	pr(5, Nuv(olaris, split("sub multi")))
	pr(6, Nuv(olaris, split("sub multi ship")))
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
	//   multi ship new <name>...
	//   multi ship <name> move <x> <y>
	//   multi ship shoot <x> <y>
	//   multi mine (set|remove) <x> <y>
	//
	// 5 <nil>
	// (multi) task [ship]
	// 6 <nil>
}
