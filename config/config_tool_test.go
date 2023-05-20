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

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigTool(t *testing.T) {
	readConfigJson := func(path string) (map[string]interface{}, error) {
		return readConfig(filepath.Join(path, "config.json"), fromConfigJson)
	}
	t.Run("new config.json", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		os.Args = []string{"config", "FOO=bar"}
		err := ConfigTool(tmpDir, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": "bar",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// t.Run("write values on existing config.json", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"foo=bar"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	err = runConfigTool([]string{"bar=baz"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{
	// 		"foo": "bar",
	// 		"bar": "baz",
	// 	}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("write existing value is overridden", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"foo=bar"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	err = runConfigTool([]string{"foo=new"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{
	// 		"foo": "new",
	// 	}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("write existing key object is merged", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"foo_bar=bar"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	err = runConfigTool([]string{"foo_baz=baz"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{
	// 		"foo": map[string]interface{}{
	// 			"bar": "bar",
	// 			"baz": "baz",
	// 		},
	// 	}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("remove existing key", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"FOO=bar"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	err = runConfigTool([]string{"FOO"}, tmpDir, true)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("remove nested key object", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"FOO_BAR=bar", "FOO_BAZ=baz"}, tmpDir, false)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	err = runConfigTool([]string{"FOO_BAR"}, tmpDir, true)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{
	// 		"foo": map[string]interface{}{
	// 			"baz": "baz",
	// 		},
	// 	}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("remove non-existing key", func(t *testing.T) {
	// 	tmpDir, _ := os.MkdirTemp("", "nuv")
	// 	defer os.RemoveAll(tmpDir)

	// 	err := runConfigTool([]string{"foo"}, tmpDir, true)
	// 	if err == nil {
	// 		t.Errorf("expected error, got nil")
	// 	}

	// 	got, err := readNuvConfigFile(tmpDir)
	// 	if err != nil {
	// 		t.Errorf("error: %v", err)
	// 	}

	// 	want := map[string]interface{}{}

	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })
}
