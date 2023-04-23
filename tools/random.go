package tools

import (
	"flag"
	"fmt"
)

func RandomGenerator() {
	flag.Usage = printRandomGeneratorHelp

	// Define command line flags
	helpFlag := flag.Bool("h", false, "Print help message")
	// integerFlag := flag.Bool("i", false, "Generate a random integer")
	// stringFlag := flag.Bool("s", false, "Generate a random string")
	// uuidFlag := flag.Bool("u", false, "Generate a random uuid v4")

	// Parse command line flags
	flag.Parse()

	// Print help message if -h flag is provided
	if *helpFlag {
		flag.Usage()
		return
	}

	// Get remaining args
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

}

func printRandomGeneratorHelp() {
	fmt.Println("Usage: random [options]")
	fmt.Println("Generates random numbers, strings and uuids")
	fmt.Println(" -h  shows this help")
	fmt.Println(" -i  <max> [min] generates a random integer between min and max (default min=0)")
	fmt.Println(" -s  <len> [<characters>] generates an alphanumeric string of length <len> from the set of <characters> provided (default <characters>=a-zA-Z0-9)")
	fmt.Println(" -u  generates a random uuid v4")
}
