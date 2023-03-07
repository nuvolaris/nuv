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

import "os"

func Example_locate() {
	dir, err := locateNuvRoot(".", false)
	pr(1, err, npath(dir))
	dir, err = locateNuvRoot("olaris", false)
	pr(2, err, npath(dir))
	dir, err = locateNuvRoot(join("olaris", "sub"), false)
	pr(3, err, npath(dir))

	// Output:
	// 1 <nil> /work/olaris
	// 2 <nil> /work/olaris
	// 3 <nil> /work/olaris
}

func Example_locate_git() {
	os.RemoveAll(join(homeDir, ".nuv"))
	dir, err := locateNuvRoot("tests", false)
	pr(1, err, nhpath(dir))
	dir, err = locateNuvRoot("tests", false)
	pr(2, err, nhpath(dir))
	// Output:
	// Cloning tasks...
	// 1 <nil> /home/.nuv/olaris
	// Updating tasks...
	// 2 <nil> /home/.nuv/olaris
}
