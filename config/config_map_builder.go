package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type configMapBuilder struct {
	configJsonPath string
	nuvRootPath    string
}

func NewConfigMapBuilder() *configMapBuilder {
	return &configMapBuilder{}
}

// WithConfigJson adds a config.json file that will be read and used
// to build a ConfigMap. If there both a config.json and a nuvroot.json are
// added, the 2 configs will be merged. It assumes the input file
// is valid.
func (b *configMapBuilder) WithConfigJson(file string) *configMapBuilder {
	b.configJsonPath = file
	return b
}

// WithNuvRoot works like WithConfigJson, with the difference that
// the NuvRoot is read and only it's inner "config":{} object is parsed
// ignoring the rest of the content.
func (b *configMapBuilder) WithNuvRoot(file string) *configMapBuilder {
	b.nuvRootPath = file
	return b
}

func (b *configMapBuilder) Build() (ConfigMap, error) {
	configJsonMap, err := readConfig(b.configJsonPath, fromConfigJson)
	if err != nil {
		return nil, err
	}

	nuvRootMap, err := readConfig(b.nuvRootPath, fromNuvRoot)
	if err != nil {
		return nil, err
	}

	return mergeMaps(configJsonMap, nuvRootMap), nil
}

func readConfig(path string, read func(string) (ConfigMap, error)) (ConfigMap, error) {
	if path == "" {
		return ConfigMap{}, nil
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ConfigMap{}, nil
	}
	if err != nil {
		return nil, err
	}

	cMap, err := read(path)
	if err != nil {
		return nil, err
	}

	return cMap, nil
}

func fromConfigJson(configPath string) (ConfigMap, error) {
	data := ConfigMap{}
	json_buf, err := os.ReadFile(configPath)
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(json_buf, &data); err != nil {
		if data == nil {
			return data, err
		}
		log.Println("config.json parsed with an error", err)
	}

	return data, nil
}

func fromNuvRoot(nuvRootPath string) (ConfigMap, error) {
	data := map[string]interface{}{}
	json_buf, err := os.ReadFile(nuvRootPath)
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(json_buf, &data); err != nil {
		if data == nil {
			return data, err
		}
		log.Println("nuvroot.json parsed with an error", err)
	}

	fmt.Printf("data: %v\n", data)
	cm, ok := data["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("config key not found in nuvroot.json")
	}
	return cm, nil
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
