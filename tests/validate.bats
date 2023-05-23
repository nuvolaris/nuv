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

@test "nuv validate help" {
    run nuv -validate
    assert_line "Usage:"

    run nuv -validate -h
    assert_line "Usage:"
}

@test "nuv validate email" {
    run nuv -validate -m example@email.com
    assert_line "'example@email.com' is a valid email address."

    run nuv -validate -m example
    assert_line "'example' is NOT a valid email address."
}

@test "nuv validate number" {
    run nuv -validate -n 123
    assert_line "'123' is a valid number."

    run nuv -validate -n 123.456
    assert_line "'123.456' is a valid number."

    run nuv -validate -n abc
    assert_line "'abc' is NOT a valid number."
}

@test "nuv validate with custom regex" {
    run nuv -validate -r '^[a-z]+$' abc
    assert_line "'abc' matches the regex."

    run nuv -validate -r '^[a-z]+$' 123
    assert_line "'123' does NOT match the regex."
}

@test "nuv validate on env vars" {
    run nuv -validate -e -n TEST_ENV_VAR
    assert_line "The variable 'TEST_ENV_VAR' is not set."

    export TEST_ENV_VAR=123
    run nuv -validate -e -n TEST_ENV_VAR
    assert_line "'123' from the variable 'TEST_ENV_VAR' is a valid number."

    run nuv -validate -e -m TEST_ENV_VAR
    assert_line "'123' from the variable 'TEST_ENV_VAR' is NOT a valid email address."

    export TEST_ENV_VAR=example@gmail.com
    run nuv -validate -e -m TEST_ENV_VAR
    assert_line "'example@gmail.com' from the variable 'TEST_ENV_VAR' is a valid email address."

    export TEST_ENV_VAR=abc
    run nuv -validate -e -r '^[a-z]+$' TEST_ENV_VAR
    assert_line "'abc' from the variable 'TEST_ENV_VAR' matches the regex."

    export TEST_ENV_VAR=123
    run nuv -validate -e -r '^[a-z]+$' TEST_ENV_VAR
    assert_line "'123' from the variable 'TEST_ENV_VAR' does NOT match the regex."
}