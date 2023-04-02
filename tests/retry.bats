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

@test "nuv -retry help" {
    run nuv -retry
    assert_success
    assert_line "Usage:"

    run nuv -retry -h
    assert_success
    assert_line "Usage:"

    run nuv -retry --help
    assert_success
    assert_line "Usage:"
}

@test "nuv -retry -t 0" {
    run nuv -retry -t 0 sub failing
    assert_line "error: failure after 0 retries or 60 seconds."
    assert_failure

    run nuv -retry -t 0 -v sub failing
    assert_line "Retry Parameters: max time=60 seconds, retries=0 times"
}

@test "nuv -retry -t 1 -m 5" {
    run nuv -retry -t 1 -m 5 sub failing
    assert_line "error: failure after 1 retries or 5 seconds."
    assert_failure

    run nuv -retry -t 1 -m 5 -v sub failing
    assert_line "Retry Parameters: max time=5 seconds, retries=1 times"
}
