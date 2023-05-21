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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func ExampleConfigTool_readValue() {
	tmpDir, _ := os.MkdirTemp("", "nuv")
	defer os.RemoveAll(tmpDir)

	nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
	configPath := filepath.Join(tmpDir, "config.json")

	os.Args = []string{"config", "FOO=bar"}
	err := ConfigTool(nuvRootPath, configPath)
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Args = []string{"config", "FOO"}
	err = ConfigTool(nuvRootPath, configPath)
	if err != nil {
		fmt.Println("error:", err)
	}

	// nested key
	os.Args = []string{"config", "NESTED_VAL=val"}
	err = ConfigTool(nuvRootPath, configPath)
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Args = []string{"config", "NESTED_VAL"}
	err = ConfigTool(nuvRootPath, configPath)
	if err != nil {
		fmt.Println("error:", err)
	}
	// Output:
	// bar
	// val
}

func TestConfigTool(t *testing.T) {
	readConfigJson := func(path string) (map[string]interface{}, error) {
		return readConfig(filepath.Join(path, "config.json"), fromConfigJson)
	}
	t.Run("new config.json", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO=bar"}
		err := ConfigTool(nuvRootPath, configPath)
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

	t.Run("write values on existing config.json", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO=bar"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		os.Args = []string{"config", "BAR=baz"}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("write existing value is overridden", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO=bar"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		os.Args = []string{"config", "FOO=new"}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": "new",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("write existing key object is merged", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO_BAR=bar"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		os.Args = []string{"config", "FOO_BAZ=baz"}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "bar",
				"baz": "baz",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("write empty string to disable key", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO_BAR=bar"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		os.Args = []string{"config", "FOO_BAR=\"\""}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("remove existing key", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO=bar"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		os.Args = []string{"config", "-r", "FOO"}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("remove nested key object", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "FOO_BAR=bar", "FOO_BAZ=baz"}
		err := ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		j, _ := readConfigJson(tmpDir)
		t.Logf("config.json: %s", j)

		os.Args = []string{"config", "-r", "FOO_BAR"}
		err = ConfigTool(nuvRootPath, configPath)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{
			"foo": map[string]interface{}{
				"baz": "baz",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("remove non-existing key", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)
		nuvRootPath := filepath.Join(tmpDir, "nuvroot.json")
		configPath := filepath.Join(tmpDir, "config.json")

		os.Args = []string{"config", "-r", "FOO"}
		err := ConfigTool(nuvRootPath, configPath)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		got, err := readConfigJson(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		want := map[string]interface{}{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func Test_buildKeyValueMap(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  keyValues
		err   error
	}{
		{
			name:  "Empty string",
			input: []string{},
			want:  nil,
			err:   fmt.Errorf("no key-value pairs provided"),
		},
		{
			name:  "Single key-value pair",
			input: []string{"foo=bar"},
			want:  keyValues{"foo": "bar"},
			err:   nil,
		},
		{
			name:  "Multiple key-value pairs",
			input: []string{"foo=bar", "baz=qux"},
			want:  keyValues{"foo": "bar", "baz": "qux"},
			err:   nil,
		},
		{
			name:  "Invalid key-value pair",
			input: []string{"foo"},
			want:  nil,
			err:   fmt.Errorf("invalid key-value pair: %q", "foo"),
		},
		{
			name:  "Invalid key-value pair",
			input: []string{"foo="},
			want:  nil,
			err:   fmt.Errorf("invalid key-value pair: %q", "foo="),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildInputKVMap(tc.input)

			if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
