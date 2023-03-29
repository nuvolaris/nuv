# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    export NO_COLOR=1
}

@test "help" {
    run nuv sub opts
    # just one as it is a cat of a message
    assert_line "Usage:"
    run nuv sub opts -h
    assert_line "Usage:"
    run nuv sub opts --help
    assert_line "Usage:"
    # do not check the actual version but ensure the output is not the help test
    run nuv sub opts --version
    refute_output "Usage:"
}

@test "cmd" {
    run nuv sub opts hello
    assert_line "hello!"
}

@test "args" {
    run nuv sub opts args mike
    assert_line "name: mike"
    assert_line "-c: no"

    run nuv sub opts args mike miri -c
    assert_line "name: mike"
    assert_line "name: miri"
    assert_line "-c: yes"
}

@test "arg1 arg3" {
    run nuv sub opts arg1 aaa arg2 1 2 --fl=ag
    assert_line "arg1 name=('aaa') arg2 x=1 y=2 --fl=ag"
    run nuv sub opts arg3 opt1 10 20 --fa
    assert_line "arg3=true opt1=true opt2=false x=10 y=20 --fa=true --fb=false"
}

@test "errors" {
    run nuv sub opts arg1
    assert_line "Usage:"
    run nuv sub opts arg1 opt4
    assert_line "Usage:"
}
