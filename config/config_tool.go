package config

import (
	"flag"
	"fmt"
	"path/filepath"
)

func ConfigTool(nuvRootDir string, configDir string) error {
	var helpFlag bool
	var dumpFlag bool
	var removeFlag bool

	flag.Usage = printConfigToolUsage

	flag.BoolVar(&helpFlag, "h", false, "show this help")
	flag.BoolVar(&helpFlag, "help", false, "show this help")
	flag.BoolVar(&dumpFlag, "dump", false, "dump the config file")
	flag.BoolVar(&dumpFlag, "d", false, "dump the config file")
	flag.BoolVar(&removeFlag, "remove", false, "remove config values")
	flag.BoolVar(&removeFlag, "r", false, "remove config values")

	flag.Parse()

	if helpFlag {
		flag.Usage()
		return nil
	}

	configMap, err := NewConfigMapBuilder().
		WithConfigJson(filepath.Join(configDir, "config.json")).
		WithNuvRoot(filepath.Join(nuvRootDir, "nuvroot.json")).
		Build()

	if err != nil {
		return err
	}

	if dumpFlag {
		dumped := configMap.Flatten()
		for k, v := range dumped {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	}

	// Get the input string from the remaining command line arguments
	input := flag.Args()

	if len(input) == 0 {
		flag.Usage()
		return nil
	}

	return runConfigTool(input, configMap, removeFlag)
}

func runConfigTool(input []string, configMap ConfigMap, removeFlag bool) error {
	return nil
}

func printConfigToolUsage() {
	fmt.Print(`Usage:
nuv -config [options] [KEY | KEY=VALUE [KEY=VALUE ...]]

Set config values passed as key-value pairs. 
If a single key is passed (without '='), its value is read from config.json, if it exists.

If you want to override a value, pass KEY="". This can be used to disable values in nuvroot.json.
Removing values from nuvroot.json is not supported, disable them instead.

-h, --help    	show this help
-r, --remove    remove config values by passing keys
-d, --dump    	dump the configs
`)
}
