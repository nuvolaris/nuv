package tools

import (
	"flag"
	"fmt"
	"os"
)

func pluginTool() error {

	flag := flag.NewFlagSet("plugin", flag.ExitOnError)
	flag.Usage = printPluginUsage

	err := flag.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return nil
	}

	fmt.Println("plugin tool not implemented yet")

	return nil
}

func printPluginUsage() {
	fmt.Println(`Usage: nuv -plugin repo

Install/update plugins from a repository.
The name of the repository must start with 'olaris-'.`)
}
