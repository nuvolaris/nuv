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
		return ConfigMap{}, err
	}

	nuvRootMap, err := readConfig(b.nuvRootPath, fromNuvRoot)
	if err != nil {
		return ConfigMap{}, err
	}

	return ConfigMap{
		nuvRootConfig: nuvRootMap,
		config:        configJsonMap,
	}, nil
}

func readConfig(path string, read func(string) (map[string]interface{}, error)) (map[string]interface{}, error) {
	if path == "" {
		return make(map[string]interface{}), nil
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return make(map[string]interface{}), nil
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

func fromConfigJson(configPath string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
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

func fromNuvRoot(nuvRootPath string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
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
