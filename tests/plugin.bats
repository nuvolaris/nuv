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
    export NUV_NO_LOG_PREFIX=1
}

@test "nuv prints 'Plugins:'" {
    run nuv
    assert_line 'Plugins:'
    assert_line "[LOCAL] olaris-test:"
    assert_line "grep-plg"
}

@test "nuv with grep plugin command" {
    run nuv grep-plg
    assert_line KO
    run nuv grep-plg GREP=first
    assert_line grep.txt:first
    assert_line OK
}

@test "nuv help on sub cmds plugin" {
    run nuv sub-plg
    assert_line '* opts:         opts test'
    assert_line '* simple:       simple'
}

@test "nuv exec sub simple plugin cmd" {
    run nuv sub-plg simple
    assert_line 'simple'
}

@test "nuv -plugin with wrong name" {
    run nuv -plugin https://github.com/giusdp/olari
    assert_line "error: plugin repository names must start with 'olaris-'"
    assert_failure
}

# @test "nuv -plugin with correct plugin repo" {
#     run nuv -plugin https://github.com/giusdp/olaris-test.git
#     assert_line "Plugins:"
#     assert_line "[LOCAL] olaris-test:"
#     assert_line "grep-plg"
# }