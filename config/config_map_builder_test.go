package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigMapBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nuv-")
	if err != nil {
		t.Error("error creating tmp dir for tests")
	}
	defer os.RemoveAll(tmpDir)

	configJsonPath := createFakeConfigFile(t, "config.json", tmpDir, `
	{
		"key": "value",
		"nested": {
			"key": 123
		}
	}`)
	nuvRootPath := createFakeConfigFile(t, "nuvroot.json", tmpDir, `
	{
		"version": "0.3.0",
		"config": {
			"nuvroot": "value",
			"another": {
				"key": 123
			}
		}
	}`)

	testCases := []struct {
		name       string
		configJson string
		nuvRoot    string
		want       ConfigMap
		err        error
	}{
		{
			name:       "should return empty configmap when no files are added",
			configJson: "",
			nuvRoot:    "",
			want:       ConfigMap{},
			err:        nil,
		},
		{
			name:       "should return configmap with values when a valid config.json is added",
			configJson: configJsonPath,
			nuvRoot:    "",
			want: ConfigMap{
				"key": "value",
				"nested": map[string]interface{}{
					"key": 123.0,
				},
			},
		},
		{
			name:       "should return configmap with values when a valid nuvroot.json is added",
			configJson: "",
			nuvRoot:    nuvRootPath,
			want: ConfigMap{
				"nuvroot": "value",
				"another": map[string]interface{}{
					"key": 123.0,
				},
			},
		},
		{
			name:       "should return configmap with values when both config.json and nuvroot.json are added",
			configJson: configJsonPath,
			nuvRoot:    nuvRootPath,
			want: ConfigMap{
				"key": "value",
				"nested": map[string]interface{}{
					"key": 123.0,
				},
				"nuvroot": "value",
				"another": map[string]interface{}{
					"key": 123.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got, err := NewConfigMapBuilder().
				WithConfigJson(tc.configJson).
				WithNuvRoot(tc.nuvRoot).
				Build()

			// if we expect an error but got none
			if tc.err != nil && err == nil {
				t.Errorf("want error %e, got %e", tc.err, err)
			}

			// if we expect no error but got one
			if tc.err == nil && err != nil {
				t.Errorf("want no error, but got %e", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func createFakeConfigFile(t *testing.T, name, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)

	if err != nil {
		t.Errorf("failed to create fake config.json: %e", err)
	}

	return path
}

func Test_mergeMaps(t *testing.T) {
	testCases := []struct {
		name     string
		m1       map[string]interface{}
		m2       map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "m1 empty",
			m1:   map[string]interface{}{},
			m2: map[string]interface{}{
				"test": "test",
			},
			expected: map[string]interface{}{
				"test": "test",
			},
		},
		{
			name: "m2 empty",
			m1: map[string]interface{}{
				"test": "test",
			},
			m2: map[string]interface{}{},
			expected: map[string]interface{}{
				"test": "test",
			},
		},
		{
			name: "m1 and m2 not empty",
			m1: map[string]interface{}{
				"test": "test",
			},
			m2: map[string]interface{}{
				"test2": "test2",
			},
			expected: map[string]interface{}{
				"test":  "test",
				"test2": "test2",
			},
		},
		{
			name: "m1 and m2 not empty with same key",
			m1: map[string]interface{}{
				"test": "test",
			},
			m2: map[string]interface{}{
				"test": "test2",
			},
			expected: map[string]interface{}{
				"test": "test2",
			},
		},
		{
			name: "m1 and m2 not empty with same key and nested map",
			m1: map[string]interface{}{
				"test": map[string]interface{}{
					"test": "test",
				},
			},
			m2: map[string]interface{}{
				"test": map[string]interface{}{
					"test2": "test2",
				},
			},
			expected: map[string]interface{}{
				"test": map[string]interface{}{
					"test":  "test",
					"test2": "test2",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mergeMaps(tc.m1, tc.m2)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
