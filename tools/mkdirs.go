package tools

import (
	"flag"
	"fmt"
	"os"
)

func Mkdirs() error {
	// Define command line flags
	helpFlag := flag.Bool("h", false, "Print help message")
	parentFlag := flag.Bool("p", false, "Create parent directories")

	// Parse command line flags
	flag.Parse()

	// Print help message if -h flag is provided
	if *helpFlag {
		printHelp()
		return nil
	}

	// Get the list of directories to create from the remaining command line arguments
	dirs := flag.Args()
	if len(dirs) == 0 {
		printHelp()
	}

	// Create each directory, with or without parent directories
	for _, dir := range dirs {
		if *parentFlag {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		} else {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printHelp() {
	fmt.Println("Usage: mkdir [-h] [-p] DIRECTORY...")
	fmt.Println("Create one or more directories.")
}
