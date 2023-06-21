package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func helpPlugins() {
	fmt.Println()
	fmt.Println("Plugins:")
}

func findPluginTask(startDir string, plg string) error {
	findTaskFunc := func(folder string) bool {
		_, err := validateTaskName(folder, plg)
		return err == nil
	}

	pluginFolder, found, err := visitPluginsFolders(startDir, findTaskFunc)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("plugin task '%s' not found", plg)
	}

	return os.Setenv("NUV_ROOT", pluginFolder)
}

func visitPluginsFolders(startDir string, predicate func(string) bool) (string, bool, error) {
	// Search in directory (startDir/olaris-*)
	dir := filepath.Join(startDir, "olaris-*")
	olarisFolders, err := filepath.Glob(dir)
	if err != nil {
		return "", false, err
	}

	for _, folder := range olarisFolders {
		if predicate(folder) {
			return folder, true, nil
		}
	}

	// Search in ~/.nuv/olaris-*
	nuvHome, err := homedir.Expand("~/.nuv")
	if err != nil {
		return "", false, err
	}
	olarisNuvFolders, err := filepath.Glob(filepath.Join(nuvHome, "olaris-*"))
	if err != nil {
		return "", false, err
	}

	for _, folder := range olarisNuvFolders {
		if predicate(folder) {
			return folder, true, nil
		}
	}

	return "", false, nil
}
