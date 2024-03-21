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

@test "-deploy with help flag" {
    run nuv -deploy -h
    assert_line --partial "Usage:"
}

@test "-deploy with missing packages folder" {
    run nuv -deploy
    assert_line --partial "no 'packages' folder found in the current directory"
    assert_failure
}

@test "-deploy with single flag with root sfa" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s hello.js -d 
    assert_line --partial "*** hello.js"
    assert_line --partial "Would run: nuv action update hello packages/hello.js"
    assert_success
}

@test "-deploy with single flag with packaged sfa" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s subfolder/hello.py -d testdata 
    assert_line --partial "*** hello.py"
    assert_line --partial "Would run: nuv package update subfolder"
    assert_line --partial "Would run: nuv action update subfolder/hello packages/subfolder/hello.py"
    assert_success
}

@test "-deploy with single flag with unsupported folder" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s subfolder -d 
    assert_line --partial "*** subfolder"
    assert_line --partial "action packages/subfolder is a directory but does not contain a supported main file"
    assert_failure
}

@test "-deploy with single flag with supported folder" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s okfolder -d 
    assert_line --partial "*** okfolder"
    assert_line --partial "Would run: nuv package update okfolder"
    assert_line --partial "Would run: nuv action update okfolder/index packages/okfolder/index.js"
    assert_success
}

@test "-deploy with single flag with unsupported MFA" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s okfolder/badmfa -d 
    assert_line --partial "*** badmfa"
    assert_line --partial "action packages/okfolder/badmfa is a directory but does not contain a supported main file"
    assert_failure
}

@test "-deploy with single flag with MFA" {
    export NUV_PWD=$(pwd)/testdata
    run nuv -deploy -s okfolder/okmfa -d 
    assert_line --partial "*** okmfa"
    assert_line --partial "Would run: nuv ide util action A=okfolder/okmfa"
    assert_line --partial "Would run: nuv package update okfolder"
    assert_line --partial "Would run: nuv action update okfolder/okmfa packages/okfolder/okmfa.zip"
    assert_success
}

@test "-deploy with scan works" {
    export NUV_PWD=$(pwd)/testdata/example
    run nuv -deploy -d 
    assert_line --partial ">>> Scan:"
    assert_line --partial "Would run: nuv ide util zip A=zipped/mfa"
    assert_line --partial "Would run: nuv ide util action A=sub/mfa"
    assert_line --partial "> packages/sub/index.js"
    assert_line --partial ">>> Deploying:"
    assert_line --partial "Would run: nuv package update zipped"
    assert_line --partial "Would run: nuv package update sub"
    assert_line --partial "Would run: nuv action update zipped/mfa packages/zipped/mfa.zip"
    assert_line --partial "Would run: nuv action update sub/mfa packages/sub/mfa.zip"
    assert_line --partial "Would run: nuv action update sub/index packages/sub/index.js "
    assert_success
}