package tools

import (
	"fmt"

	"github.com/nuvolaris/someutils"
)

var Utils = []string{
	"basename", "cat", "cp", "dirname",
	"grep", "gunzip", "gzip", "head",
	"ls", "mv", "pwd", "rm", "sleep",
	"tail", "tar", "tee", "touch", "tr",
	"unzip", "wc", "which", "zip",
}

func IsUtil(name string) bool {
	for _, s := range Utils {
		if s == name {
			return true
		}
	}
	return false
}

func RunUtil(name string, args []string) (int, error) {
	if !IsUtil(name) {
		return 1, fmt.Errorf("command %s not found", name)
	}
	full := append([]string{name}, args...)
	fmt.Println(name, full)
	err, code := someutils.Call(name, full)
	return code, err
}
