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

@test "-deploy without args" {
    run nuv -deploy
    assert_line --partial "Usage:"
}

@test "-deploy with help flag" {
    run nuv -deploy -h
    assert_line --partial "Usage:"
}

@test "-deploy with missing packages folder" {
    run nuv -deploy .
    assert_line --partial "no 'packages' folder found in the given directory"
    assert_failure
}

@test "-deploy with single flag with root sfa" {
    run nuv -deploy -s hello.js -d testdata 
    assert_line --partial "*** hello.js"
    assert_line --partial "Would run: nuv action update hello testdata/packages/hello.js"
    assert_success
}

@test "-deploy with single flag with packaged sfa" {
    run nuv -deploy -s subfolder/hello.py -d testdata 
    assert_line --partial "*** hello.py"
    assert_line --partial "Would run: nuv package update subfolder"
    assert_line --partial "Would run: nuv action update subfolder/hello testdata/packages/subfolder/hello.py"
    assert_success
}

@test "-deploy with single flag with unsupported folder" {
    run nuv -deploy -s subfolder -d testdata 
    assert_line --partial "*** subfolder"
    assert_line --partial "action packages/subfolder is a directory but does not contain a supported main file"
    assert_failure
}

@test "-deploy with single flag with supported folder" {
    run nuv -deploy -s okfolder -d testdata 
    assert_line --partial "*** okfolder"
    assert_line --partial "Would run: nuv package update okfolder"
    assert_line --partial "Would run: nuv action update okfolder/index testdata/packages/okfolder/index.js"
    assert_success
}

@test "-deploy with single flag with unsupported MFA" {
    run nuv -deploy -s okfolder/badmfa -d testdata 
    assert_line --partial "*** badmfa"
    assert_line --partial "action packages/okfolder/badmfa is a directory but does not contain a supported main file"
    assert_failure
}

@test "-deploy with single flag with MFA" {
    run nuv -deploy -s okfolder/okmfa -d testdata 
    assert_line --partial "*** okmfa"
    assert_line --partial "Would run: nuv ide util action A=okfolder/okmfa"
    assert_line --partial "Would run: nuv package update okfolder"
    assert_line --partial "Would run: nuv action update okfolder/okmfa testdata/packages/okfolder/okmfa.zip"
    assert_success
}