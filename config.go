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
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type keyValues map[string]string

func (kv *keyValues) String() string {
	return fmt.Sprintf("%v", *kv)
}

func (kv *keyValues) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid key-value pair: %q", value)
	}
	key := parts[0]
	val := parts[1]

	if key == "" || val == "" {
		return fmt.Errorf("invalid key-value pair: %q", value)
	}

	if *kv == nil {
		*kv = make(keyValues)
	}
	(*kv)[key] = val
	return nil
}

func printConfigToolUsage() {
	fmt.Print(`Usage:
nuv -config [-h|--help] KEY=VALUE [KEY=VALUE ...]

Set config values passed as key-value pairs.

-h, --help    show this help
-d, --dump    dump the configs
`)
}

func ConfigTool() error {
	var helpFlag bool
	var dumpFlag bool

	flag.Usage = printConfigToolUsage

	flag.BoolVar(&helpFlag, "h", false, "show this help")
	flag.BoolVar(&helpFlag, "help", false, "show this help")
	flag.BoolVar(&dumpFlag, "dump", false, "dump the config file")
	flag.BoolVar(&dumpFlag, "d", false, "dump the config file")

	flag.Parse()

	if helpFlag {
		flag.Usage()
		return nil
	}

	if dumpFlag {
		return dumpAll()
	}

	// Get the input string from the remaining command line arguments
	input := flag.Args()

	if len(input) == 0 {
		flag.Usage()
		return nil
	}

	home, err := homedir.Expand("~/.nuv")
	if err != nil {
		return err
	}

	return runConfigTool(input, home)
}

func runConfigTool(input []string, dir string) error {
	config := make(map[string]interface{})
	// Check if the config file exists
	if exists(dir, CONFIGFILE) {
		// If it exists, load it
		configFromFile, err := readNuvConfigFile(dir)
		if err != nil {
			return err
		}
		config = configFromFile
	}

	kv, err := buildKeyValueMap(input)
	if err != nil {
		return err
	}
	newConfig, err := buildConfigObject(kv)
	if err != nil {
		return err
	}

	// Merge the input config into the existing config
	// NOTE: input keys are merged with existing config.json keys
	// with priority given to the new keys (in case of conflicts)
	config = mergeMaps(config, newConfig)

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write the config file
	err = os.WriteFile(joinpath(dir, CONFIGFILE), configJSON, 0644)
	if err != nil {
		warn("failed to write config.json")
		return err
	}
	return err
}

func dumpAll() error {
	home, err := homedir.Expand("~/.nuv")
	if err != nil {
		return err
	}
	nuvRootDir, err := getNuvRoot()
	if err != nil {
		return err
	}

	nuvRoot, err := readNuvRootFile(nuvRootDir)
	if err != nil {
		return err
	}
	nuvRootConfigEnv := rebuildConfigEnvVars(nuvRoot.Config)

	configs, err := readNuvConfigFile(home)
	if err != nil {
		return err
	}
	configEnv := rebuildConfigEnvVars(configs)

	for k, v := range nuvRootConfigEnv {
		if _, ok := configEnv[k]; !ok {
			key := strings.ToUpper(k)
			fmt.Printf("%s=%s\n", key, v)
		}
	}
	for k, v := range configEnv {
		key := strings.ToUpper(k)
		fmt.Printf("%s=%s\n", key, v)
	}

	return nil
}

func buildConfigObject(kv keyValues) (map[string]interface{}, error) {

	// KEYs form the keys in the config json file, delimiting inner objects by the underscore character.
	// For example, the key "foo_bar" will be interpreted as the key "bar" in the object "foo".
	// If 2 keys have the same prefix, they will be merged into the same object.

	configToSave := make(map[string]interface{})

	for k, value := range kv {
		// Key parsing
		keys, err := parseKey(strings.ToLower(k))
		if err != nil {
			return nil, err
		}
		lastIndex := len(keys) - 1

		currentMap := configToSave
		for i, subKey := range keys {
			// If we are at the last key, set the value
			if i == lastIndex {
				v, err := parseValue(value)
				if err != nil {
					return nil, err
				}

				// If the key already exists, return an error
				if _, ok := currentMap[subKey]; ok {
					return nil, fmt.Errorf("invalid key: '%s' - '%s' is already being used for a value", k, subKey)
				}

				currentMap[subKey] = v
			} else {
				// If the sub-map doesn't exist, create it
				if _, ok := currentMap[subKey]; !ok {
					currentMap[subKey] = make(map[string]interface{})
				}
				// Update the current map to the sub-map
				m, ok := currentMap[subKey].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid key: '%s' - '%s' is already being used for a value", k, subKey)
				}
				currentMap = m
			}
		}
	}

	return configToSave, nil
}

func parseKey(key string) ([]string, error) {
	parts := strings.Split(key, "_")
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid key: %s", key)
		}
	}
	return parts, nil
}

/*
VALUEs are parsed in the following way:

  - try to parse as a jsos first, and if it is a json, store as a json
  - then try to parse as a number, and if it is a (float) number store as a number
  - then try to parse as true or false and store as a boolean
  - then check if it's null and store as a null
  - otherwise store as a string
*/
func parseValue(value string) (interface{}, error) {
	// Try to parse as json
	var jsonValue interface{}
	if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
		return jsonValue, nil
	}

	// Try to parse as a integer with strconv
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue, nil
	}

	// Try to parse as a float with strconv
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue, nil
	}

	// Try to parse as a boolean
	if value == "true" || value == "false" {
		return value == "true", nil
	}

	// Try to parse as null
	if value == "null" {
		return nil, nil
	}

	// Otherwise, return the string
	return value, nil
}

func buildKeyValueMap(pairs []string) (keyValues, error) {
	var kv keyValues

	if len(pairs) == 0 {
		return nil, fmt.Errorf("no key-value pairs provided")
	}

	for _, pair := range pairs {

		if err := kv.Set(pair); err != nil {
			return nil, err
		}
	}
	return kv, nil
}

func applyAllConfigEnvVars(nuvRootDir, configJsonDir string) error {
	if err := applyNuvRootConfigEnvVars(nuvRootDir); err != nil {
		return err
	}

	if err := applyConfigJsonEnvVars(configJsonDir); err != nil {
		return err
	}

	return nil
}

func applyNuvRootConfigEnvVars(nuvRootDir string) error {
	trace("Applying env vars from nuvroot.json...")
	nuvRoot, err := readNuvRootFile(nuvRootDir)
	if err != nil {
		return err
	}

	nuvRootConfigEnv := rebuildConfigEnvVars(nuvRoot.Config)

	for k, v := range nuvRootConfigEnv {
		key := strings.ToUpper(k)
		os.Setenv(key, v)
		debug("Set env var (nuvroot.json)", key, "=", v)
	}

	return nil
}

func applyConfigJsonEnvVars(configJsonDir string) error {
	trace("Applying env vars from config.json...")
	configs, err := readNuvConfigFile(configJsonDir)
	if err != nil {
		return err
	}
	configEnv := rebuildConfigEnvVars(configs)

	for k, v := range configEnv {
		key := strings.ToUpper(k)
		os.Setenv(key, v)
		debug("Set env var (config.json)", key, "=", v)
	}

	return nil
}

func rebuildConfigEnvVars(config map[string]interface{}) map[string]string {
	outputMap := make(map[string]string)
	traverse(config, "", outputMap)
	return outputMap
}

func traverse(inputMap map[string]interface{}, prefix string, outputMap map[string]string) {
	for k, v := range inputMap {
		if m, ok := v.(map[string]interface{}); ok {
			if prefix == "" {
				traverse(m, k, outputMap)
			} else {
				traverse(m, fmt.Sprintf("%s_%s", prefix, k), outputMap)
			}
		} else {
			if prefix == "" {
				outputMap[k] = fmt.Sprintf("%v", v)
			} else {
				outputMap[fmt.Sprintf("%s_%s", prefix, k)] = fmt.Sprintf("%v", v)
			}
		}
	}
}

// mergeMaps merges map2 into map1 overwriting any values in map1 with values from map2
// when there are conflicts. It returns the merged map.
func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})

	for key, value := range map1 {

		map2Value, ok := map2[key]
		// key doesn't exist in map2 so add it to the merged map
		if !ok {
			mergedMap[key] = value
			continue
		}

		// key exists in map2 but map1 value is NOT a map, so add value from map2
		mapFromMap1, ok := value.(map[string]interface{})
		if !ok {
			mergedMap[key] = map2Value
			continue
		}

		mapFromMap2, ok := map2Value.(map[string]interface{})
		// key exists in map2, map1 value IS a map but map2 value is not, so overwrite with map2
		if !ok {
			mergedMap[key] = mapFromMap2
			continue
		}

		// key exists in map2, map1 value IS a map, map2 value IS a map, so merge recursively
		mergedMap[key] = mergeMaps(mapFromMap1, mapFromMap2)
	}

	// add any keys that exist in map2 but not in map1
	for key, value := range map2 {
		if _, ok := mergedMap[key]; !ok {
			mergedMap[key] = value
		}
	}

	return mergedMap
}
