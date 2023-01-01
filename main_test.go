package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

/// test utils

// pr print args
func pr(args ...any) {
	fmt.Println(args...)
}

// as creates a string array
func as(s ...string) []string {
	return s
}

var homeDir = ""

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	homeDir, _ = filepath.Abs(wd)
	os.Exit(m.Run())
}
