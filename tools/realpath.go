package tools

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func realpathTool() error {
	flags := flag.NewFlagSet("realpath", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Println("Usage: nuv -realpath <path>")
		fmt.Println()
		fmt.Println("Options:")
		flags.PrintDefaults()
	}

	showHelp := flags.Bool("h", false, "Show help")

	// Parse command-line arguments
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *showHelp {
		flags.Usage()
		return nil
	}

	if flags.NArg() != 1 {
		flags.Usage()
		return errors.New("no path provided")
	}

	path := flags.Arg(0)

	var absPath string
	if filepath.IsLocal(path) {
		absPath = filepath.Join(os.Getenv("NUV_PWD"), path)
	} else {
		absPath = path
	}

	absPath = filepath.Clean(absPath)
	fmt.Println(absPath)

	return nil
}
