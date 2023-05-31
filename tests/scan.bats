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

@test "-scan -h" {
    run nuv -scan -h
    assert_line "Usage:"

    run nuv -scan
    assert_line "Usage:"
    assert_failure
}

@test "-scan stops if actions folder not present" {
    NUV_DIR="./olaris" run  nuv -scan nuv -js
    assert_line --partial "actions folder not found"
    assert_failure
}

# @test "-scan -js argv.js" {
#     run NUV_DIR="./testdata" nuv -scan -js testdata/js_test_argv.js
#     assert_line --partial "/testdata/actions"
#     assert_line --partial "/testdata/actions/subfolder"
#     assert_line --partial "/testdata/actions/subfolder/subsub hello.js"
#     assert_success
# }