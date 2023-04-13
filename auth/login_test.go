// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package auth

import (
	"fmt"
	"testing"

	"github.com/zalando/go-keyring"
)

// TODO:
// test with bad apihost <-- a url regex validation would be nice
// test with mock http server

func ExampleLoginCmd_noArgs() {
	err := LoginCmd([]string{})
	fmt.Println(err)
	// Output:
	// Usage:
	// nuv login <apihost> [<user>]
	// <nil>
}

func ExampleLoginCmd_askPasswordSuccess() {
	oldPwdReader := pwdReader
	pwdReader = &stubPasswordReader{
		Password:    "fakePassword",
		ReturnError: false,
	}

	err := LoginCmd([]string{"fakeApiHost", "fakeUser"})
	pwdReader = oldPwdReader
	fmt.Println(err)
	// Output:
	// Enter Password: <nil>
}

func TestLoginCmd(t *testing.T) {
	t.Run("error: returns error when empty password", func(t *testing.T) {
		oldPwdReader := pwdReader
		pwdReader = &stubPasswordReader{
			Password:    "",
			ReturnError: false,
		}

		err := LoginCmd([]string{"fakeApiHost", "fakeUser"})
		pwdReader = oldPwdReader
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "password is empty" {
			t.Errorf("Expected error to be 'password is empty', got %s", err.Error())
		}
	})
}

func Example_storeCredentials() {
	keyring.MockInit()

	fakeCreds := make(map[string]string)
	fakeCreds["AUTH"] = "fakeValue1"
	fakeCreds["REDIS_URL"] = "fakeValue2"
	fakeCreds["MONGODB"] = "fakeValue3"

	err := storeCredentials(fakeCreds)
	fmt.Println(err)
	for k := range fakeCreds {
		cred, err := keyring.Get(nuvSecretServiceName, k)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(k + ": " + cred)
	}

	// Output:
	// <nil>
	// AUTH: fakeValue1
	// REDIS_URL: fakeValue2
	// MONGODB: fakeValue3
}
