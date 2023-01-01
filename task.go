package main

import (
	"fmt"
	"os"
	"path"

	"github.com/nuvolaris/task/cmd/taskmain/v3"
)

var taskDryRun = false

func Task(args ...string) {
	if taskDryRun {
		cur, _ := os.Getwd()
		dir := path.Base(cur)
		fmt.Printf("(%s) task %v\n", dir, args)
	} else {
		taskmain.Task(args)
	}
}
