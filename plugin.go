package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nuvolaris/task/v3"
)

func printInstalledPluginsMessage(localDir string) error {
	plgs, err := newPlugins(localDir)
	if err != nil {
		return err
	}
	plgs.print()
	return nil
}

func findPluginTask(startDir string, plg string) error {
	trace("findPluginTask", startDir, plg)
	// findTaskFunc := func(folder string) bool {
	// 	_, err := validateTaskName(folder, plg)
	// 	return err == nil
	// }

	// pluginFolder, found, err := visitPluginsFolders(startDir, findTaskFunc)
	// if err != nil {
	// 	return err
	// }
	// if !found {
	// 	return fmt.Errorf("plugin task '%s' not found", plg)
	// }

	// return os.Setenv("NUV_ROOT", pluginFolder)
	return nil
}

type plugins struct {
	local []string
	nuv   []string
}

func newPlugins(localDir string) (*plugins, error) {
	localOlarisFolders := make([]string, 0)
	nuvOlarisFolders := make([]string, 0)

	// Search in directory (localDir/olaris-*)
	dir := filepath.Join(localDir, "olaris-*")
	olarisFolders, err := filepath.Glob(dir)
	if err != nil {
		return nil, err
	}

	localOlarisFolders = append(localOlarisFolders, olarisFolders...)

	// Search in ~/.nuv/olaris-*
	nuvHome, err := homedir.Expand("~/.nuv")
	if err != nil {
		return nil, err
	}

	olarisNuvFolders, err := filepath.Glob(filepath.Join(nuvHome, "olaris-*"))
	if err != nil {
		return nil, err
	}

	nuvOlarisFolders = append(nuvOlarisFolders, olarisNuvFolders...)

	return &plugins{
		local: localOlarisFolders,
		nuv:   nuvOlarisFolders,
	}, nil
}

func (p *plugins) print() {
	if len(p.local) == 0 && len(p.nuv) == 0 {
		fmt.Println("No plugins installed. Use 'nuv -plugin' to add new ones.")
	}

	fmt.Println("Plugins:")
	if len(p.local) > 0 {
		for _, plg := range p.local {
			fmt.Printf("[LOCAL] %s:\n", filepath.Base(plg))
			printTaskHelp(plg)
		}
	}

	if len(p.nuv) > 0 {
		for _, plg := range p.nuv {
			fmt.Printf("  %s:\n", filepath.Base(plg))
			printTaskHelp(plg)
		}
	}
}

func printTaskHelp(path string) {
	dir := path
	entrypoint := NUVFILE
	e := task.Executor{
		Dir:        dir,
		Entrypoint: entrypoint,
		Summary:    true,
		Color:      true,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	e.ListTaskNames(true)
}
