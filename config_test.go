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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_runConfigTool(t *testing.T) {
	t.Run("new config.json", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		err := runConfigTool([]string{"foo=bar"}, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readNuvConfigFile(tmpDir)
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

	t.Run("existing config.json", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		err := runConfigTool([]string{"foo=bar"}, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		err = runConfigTool([]string{"bar=baz"}, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readNuvConfigFile(tmpDir)
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

	t.Run("existing config.json is overridden", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		err := runConfigTool([]string{"foo=bar"}, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		err = runConfigTool([]string{"foo=new"}, tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		got, err := readNuvConfigFile(tmpDir)
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

}

func Test_buildConfigObject(t *testing.T) {
	testCases := []struct {
		name  string
		input keyValues
		want  map[string]interface{}
		err   error
	}{
		{
			name: "Simple Key",
			input: keyValues{
				"foo": "bar",
			},
			want: map[string]interface{}{
				"foo": "bar",
			},
			err: nil,
		},
		{
			name: "Complex Key",
			input: keyValues{
				"foo_bar": "baz",
			},
			want: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			},
			err: nil,
		},
		{
			name: "Multiple Keys",
			input: keyValues{
				"foo_bar": "baz",
				"foo_baz": "bar",
			},
			want: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
					"baz": "bar",
				},
			},
			err: nil,
		},
		{
			name: "Duplicate key",
			input: keyValues{
				"foo":     "bar",
				"foo_bar": "baz",
			},
			want: nil,
			err:  fmt.Errorf("invalid key: '%s' - '%s' is already being used for a value", "foo_bar", "foo"),
		},
		{
			name: "Duplicate key 2",
			input: keyValues{
				"foo_bar_baz": "bar",
				"foo_bar":     "baz",
			},
			want: nil,
			err:  fmt.Errorf("invalid key: '%s' - '%s' is already being used for a value", "foo_bar", "bar"),
		},
		{
			name: "Invalid key",
			input: keyValues{
				"foo_bar_": "baz",
			},
			want: nil,
			err:  fmt.Errorf("invalid key: %s", "foo_bar_"),
		},
		{
			name: "Invalid key 2",
			input: keyValues{
				"_foo_bar": "baz",
			},
			want: nil,
			err:  fmt.Errorf("invalid key: %s", "_foo_bar"),
		},
		{
			name: "Invalid key 3",
			input: keyValues{
				"foo__bar": "baz",
			},
			want: nil,
			err:  fmt.Errorf("invalid key: %s", "foo__bar"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildConfigObject(tc.input)

			if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}

		})
	}
}

func Test_parseKey(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
		err   error
	}{
		{
			name:  "Simple Key",
			input: "foo",
			want:  []string{"foo"},
			err:   nil,
		},
		{
			name:  "Complex Key",
			input: "foo_bar",
			want:  []string{"foo", "bar"},
			err:   nil,
		},
		{
			name:  "Complex Key 2",
			input: "foo_bar_baz",
			want:  []string{"foo", "bar", "baz"},
			err:   nil,
		},
		{
			name:  "Invalid Key",
			input: "foo_bar_baz_",
			want:  nil,
			err:   fmt.Errorf("invalid key: %s", "foo_bar_baz_"),
		},
		{
			name:  "Invalid Key 2",
			input: "_foo_bar_baz",
			want:  nil,
			err:   fmt.Errorf("invalid key: %s", "_foo_bar_baz"),
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseKey(tc.input)

			if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func Test_parseValue(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  interface{}
		err   error
	}{
		{
			name:  "String",
			input: "foo",
			want:  "foo",
			err:   nil,
		},
		{
			name:  "Complex String",
			input: "Another foo bar",
			want:  "Another foo bar",
			err:   nil,
		},
		{
			name:  "Number",
			input: "123.456",
			want:  123.456,
			err:   nil,
		},
		{
			name:  "Boolean True",
			input: "true",
			want:  true,
			err:   nil,
		},
		{
			name:  "Boolean False",
			input: "false",
			want:  false,
			err:   nil,
		},
		{
			name:  "Null",
			input: "null",
			want:  nil,
			err:   nil,
		},
		{
			name:  "JSON",
			input: `{"foo": "bar"}`,
			want:  map[string]interface{}{"foo": "bar"},
			err:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseValue(tc.input)

			if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
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
			got, err := buildKeyValueMap(tc.input)

			if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func Test_rebuildConfigEnvVars(t *testing.T) {
	testCases := []struct {
		name  string
		input map[string]interface{}
		want  map[string]string
	}{
		{
			name:  "Empty",
			input: map[string]interface{}{},
			want:  map[string]string{},
		},
		{
			name:  "Single env var",
			input: map[string]interface{}{"foo": "bar"},
			want:  map[string]string{"foo": "bar"},
		},
		{
			name:  "Multiple env vars",
			input: map[string]interface{}{"foo": "bar", "baz": "qux"},
			want:  map[string]string{"foo": "bar", "baz": "qux"},
		},
		{
			name:  "Nested env vars",
			input: map[string]interface{}{"foo": map[string]interface{}{"bar": "baz"}},
			want:  map[string]string{"foo_bar": "baz"},
		},
		{
			name:  "Nested env vars 2",
			input: map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "qux"}}},
			want:  map[string]string{"foo_bar_baz": "qux"},
		},
		{
			name:  "Multiple nested env vars",
			input: map[string]interface{}{"foo": map[string]interface{}{"bar": "baz", "qux": "quux"}},
			want:  map[string]string{"foo_bar": "baz", "foo_qux": "quux"},
		},
		{
			name:  "Multiple nested env vars 2",
			input: map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "qux"}, "quux": "corge"}},
			want:  map[string]string{"foo_bar_baz": "qux", "foo_quux": "corge"},
		},
		{
			name:  "Multiple nested env vars 3",
			input: map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "qux"}, "quux": map[string]interface{}{"corge": "grault"}}},
			want:  map[string]string{"foo_bar_baz": "qux", "foo_quux_corge": "grault"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := rebuildConfigEnvVars(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}

}

func Test_applyNuvRootConfigEnvVars(t *testing.T) {
	t.Run("no env vars set when nuvroot.json config not set", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		createTestNuvRootFile(t, tmpDir, true)

		err := applyNuvRootConfigEnvVars(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		// assert no env vars set
		for k := range testEnvVars {
			if os.Getenv(k) != "" {
				t.Errorf("env var %s should not be set", k)
			}
		}
	})

	t.Run("env vars from nuvroot.json config when present ", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		createTestNuvRootFile(t, tmpDir, false)

		err := applyNuvRootConfigEnvVars(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		assertEnvVar(t)

		// reset vars
		for k := range testEnvVars {
			os.Setenv(k, "")
		}
	})
}

func Test_applyConfigJsonEnvVars(t *testing.T) {

	t.Run("env vars from config.json when exists", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "nuv")
		defer os.RemoveAll(tmpDir)

		// create config.json
		createTestConfigJson(t, tmpDir, testEnvVars)

		err := applyConfigJsonEnvVars(tmpDir)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		assertEnvVar(t)

		// reset vars
		for k := range testEnvVars {
			os.Setenv(k, "")
		}
	})

}

func Test_applyAllConfigEnvVars(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "nuv")
	defer os.RemoveAll(tmpDir)

	// create nuvroot.json
	createTestNuvRootFile(t, tmpDir, false)

	// create config.json
	overrideTestEnvVar := map[string]interface{}{
		"test": map[string]interface{}{
			"env": map[string]interface{}{
				"var2": "test2-override",
			},
			"var3": "test3-override",
		},
	}

	createTestConfigJson(t, tmpDir, overrideTestEnvVar)

	err := applyAllConfigEnvVars(tmpDir, tmpDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if os.Getenv("NUV_TEST_ENV_VAR") != "test" {
		t.Errorf("env var NUV_TEST_ENV_VAR should be set to test")
	}

	if os.Getenv("TEST_ENV_VAR2") != "test2-override" {
		t.Errorf("env var TEST_ENV_VAR2 should be set to test2-override")
	}

	if os.Getenv("TEST_VAR3") != "test3-override" {
		t.Errorf("env var TEST_VAR3 should be set to test3-override")
	}

	if os.Getenv("TEST4") != "test4" {
		t.Errorf("env var TEST4 should be set to test4")
	}

	// reset vars
	for k := range testEnvVars {
		os.Setenv(k, "")
	}
}

var testExpectedEnvVars = map[string]string{
	"NUV_TEST_ENV_VAR": "test",
	"TEST_ENV_VAR2":    "test2",
	"TEST_VAR3":        "test3",
	"TEST4":            "test4",
}

var testEnvVars = map[string]interface{}{
	"nuv": map[string]interface{}{
		"test": map[string]interface{}{
			"env": map[string]interface{}{
				"var": "test",
			},
		},
	},
	"test": map[string]interface{}{
		"env": map[string]interface{}{
			"var2": "test2",
		},
		"var3": "test3",
	},
	"test4": "test4",
}

func assertEnvVar(t *testing.T) {
	t.Helper()
	for k, v := range testExpectedEnvVars {
		if os.Getenv(k) != v {
			t.Errorf("env var %s should be set to %s", k, v)
		}
	}
}

func createTestNuvRootFile(t *testing.T, dir string, onlyVersion bool) {
	t.Helper()
	path := filepath.Join(dir, "nuvroot.json")
	var nuvRootStr string
	if onlyVersion {
		nuvRootStr = "{ \"version\": \"0.3.0\" }"
	} else {
		configJson, err := json.Marshal(testEnvVars)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		nuvRootStr = fmt.Sprintf("{ \"version\": \"0.3.0\", \"config\": %s }", configJson)
	}
	os.WriteFile(path, []byte(nuvRootStr), 0644)
}

func createTestConfigJson(t *testing.T, dir string, envVars map[string]interface{}) {
	t.Helper()
	configJson, err := json.Marshal(envVars)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	os.WriteFile(filepath.Join(dir, "config.json"), []byte(configJson), 0644)
}
