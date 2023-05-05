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

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIndex(t *testing.T) {
	_ = os.Chdir(workDir)
	olaris, _ := filepath.Abs(joinpath("tests", "olaris"))
	handler := webFileServerHandler(joinpath(olaris, WebDir))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		r    *http.Request
	}{
		{name: "1: testing get", r: newreq("GET", ts.URL+"/", nil)},
		// {name: "2: testing post", r: newreq("POST", ts.URL+"/", nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.r)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
			resp.Body.Close()
		})
	}
}

func TestGetTask(t *testing.T) {
	_ = os.Chdir(joinpath(workDir, "tests"))

	t.Run("/api/nuv?test invokes test task", func(t *testing.T) {
		request := newTaskRequest("test")
		response := httptest.NewRecorder()

		want := NuvOutput{
			Stdout: "Dry run: task test",
			Stderr: "",
			Status: 0,
		}

		nuvTaskServer(response, request)

		assertResponseBody(t, response.Body.Bytes(), want)
	})
}

func newTaskRequest(task string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/nuv?"+task, nil)
	return req
}

func assertResponseBody(t testing.TB, body []byte, want NuvOutput) {
	got := NuvOutput{}
	err := json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}
	t.Helper()
	if got.Status != want.Status {
		t.Errorf("expected status %v; got %v", want.Status, got.Status)
	}
	if got.Stdout != want.Stdout {
		t.Errorf("expected stdout %v; got %v", want.Stdout, got.Stdout)
	}
	if got.Stderr != want.Stderr {
		t.Errorf("expected stderr %v; got %v", want.Stderr, got.Stderr)
	}
}
