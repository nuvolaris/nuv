package config

import (
	"fmt"
	"strings"
)

// A ConfigMap is a map where the keys are in the form of: A_KEY_WITH_UNDERSCORES.
// The map splits the key by the underscores and creates a nested map that
// represents the key. For example, the key "A_KEY_WITH_UNDERSCORES" would be
// represented as:
//
//	{
//		"a": {
//			"key": {
//				"with": {
//					"underscores": "value",
//				},
//			},
//		},
//	}
//
// To interact with the ConfigMap, use the Insert, Get, Update, and Delete by passing
// keys in the form above.
type ConfigMap map[string]interface{}

// Insert inserts a key and value into the ConfigMap. If the key already exists,
// the value is overwritten. The expected key format is A_KEY_WITH_UNDERSCORES.
func (c *ConfigMap) Insert(key string, value interface{}) error {

	_, err := parseKey(strings.ToLower(key))
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigMap) Get(key string) interface{} {
	return nil
}

func (c *ConfigMap) Update(key string, value interface{}) bool {
	return false
}

func (c *ConfigMap) Delete(key string) bool {
	return false
}

func (c *ConfigMap) DumpAll() {

}

func (c *ConfigMap) Keys() []string {
	return nil
}

func (c *ConfigMap) flatten() map[string]string {
	return map[string]string{}
}

////

func parseKey(key string) ([]string, error) {
	parts := strings.Split(key, "_")
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid key: %s", key)
		}
	}
	return parts, nil
}
