package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInsert(t *testing.T) {

	testCases := []struct {
		name        string
		startingMap map[string]interface{}
		key         string
		value       interface{}
		expected    map[string]interface{}
	}{
		{
			name:        "empty map",
			startingMap: map[string]interface{}{},
			key:         "KEY",
			value:       "value",
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "map with key",
			startingMap: map[string]interface{}{
				"key": "value",
			},
			key:   "KEY2",
			value: "value2",
			expected: map[string]interface{}{
				"key":  "value",
				"key2": "value2",
			},
		},
		{
			name: "map with nested key",
			startingMap: map[string]interface{}{
				"key": map[string]interface{}{
					"key": "value",
				},
			},
			key:   "KEY_OTHER",
			value: "value2",
			expected: map[string]interface{}{
				"key": map[string]interface{}{
					"key":   "value",
					"other": "value2",
				},
			},
		},
		{
			name: "overwrite value",
			startingMap: map[string]interface{}{
				"key": "value",
			},
			key:   "KEY",
			value: "value2",
			expected: map[string]interface{}{
				"key": "value2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := ConfigMap(tc.startingMap)
			c.Insert(tc.key, tc.value)
			if !reflect.DeepEqual(c, tc.expected) {
				t.Errorf("expected: %v, got: %v", tc.expected, c)
			}
		})
	}
}

// //
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
