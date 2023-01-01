package main

import "os"

func ExampleNuvSubdirs() {
	os.Chdir(homeDir)
	pr(1, Nuv("notexistent", as()))
	os.Chdir(homeDir)
	pr(2, Nuv("tests", as()))
	os.Chdir(homeDir)
	pr(3, Nuv("tests", as("sub")))
	os.Chdir(homeDir)
	pr(4, Nuv("tests", as("sub", "subsub")))
	// Output:
	// -
}

func ExampleNuvCmd() {
	os.Chdir(homeDir)
	pr(1, Nuv("tests", as("hello")))
	// Output:
	// -
}
