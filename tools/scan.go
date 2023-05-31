package tools

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// The scan tool visits a folder and creates an action plan to execute the given cmd (called cmdPlan).
// The cmd is the given nuv command to run, and the args is an array of arrays.
// Each entry is one execution instance, which holds the folder path (first arg) and files names (the rest).
// The plan is then executed by running the cmd once for each entry of args (so args.length times)
// in the form of `cmd args[i][0] args[i][1] ... args[i][n]` for each i in args.
type cmdPlan struct {
	cmd  string
	args [][]string
}

func (a *cmdPlan) setCmd(cmd string) {
	a.cmd = cmd
}
func (a *cmdPlan) appendArg(args []string) {
	a.args = append(a.args, args)
}

const actionsFolder = "actions"

func scanTool() error {
	flag := flag.NewFlagSet("scan", flag.ExitOnError)
	flag.Usage = printScanUsage

	helpFlag := flag.Bool("h", false, "show help")
	dirFlag := flag.String("d", getCurrentDir(), "directory to scan")

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *helpFlag {
		flag.Usage()
		return nil
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return errors.New("missing required nuv command")
	}

	p, err := filepath.Abs(*dirFlag)
	if err != nil {
		return err
	}

	_, err = buildActionPlan(p, args)
	if err != nil {
		return err
	}

	return nil
}

func buildActionPlan(scanPath string, cmd []string) (*cmdPlan, error) {
	// check if actions folder exists
	if err := checkActionsFolder(scanPath); err != nil {
		return nil, err
	}

	_, err := getAllDirs(scanPath)
	if err != nil {
		return nil, err
	}

	// res, err := visitScanDir(fsys, "", actionsFolder)
	// if err != nil {
	// 	return nil, err
	// }

	// plan := &cmdPlan{}
	// plan.setCmd(cmd[0])
	// for _, r := range res {
	// 	plan.appendArg(r)
	// }

	return nil, nil
}

func getAllDirs(rootDir string) ([]string, error) {
	var dirs []string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	return dirs, err
}

func getAllFiles(dir string) ([]string, error) {
	var files []string

	children, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, child := range children {
		if !child.IsDir() {
			files = append(files, child.Name())
		}
	}

	return files, err
}

func visitScanDir(fsys fs.FS, parentPath string, childPath string) ([][]string, error) {
	var result [][]string

	// fs.WalkDir()
	dirPath, err := filepath.Abs(filepath.Join(parentPath, childPath))
	if err != nil {
		return nil, err
	}
	children, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return nil, err
	}

	var dirList []string
	var fileList []string
	for _, child := range children {
		if child.IsDir() {
			dirList = append(dirList, child.Name())
		} else {
			fileList = append(fileList, child.Name())
		}
	}
	result = append(result, []string{dirPath}, fileList)

	for _, subDir := range dirList {
		subResult, err := visitScanDir(fsys, dirPath, subDir)
		if err != nil {
			return nil, err
		}

		result = append(result, subResult...)
	}

	return result, nil
}

func checkActionsFolder(path string) error {
	info, err := os.Stat(filepath.Join(path, actionsFolder))
	if os.IsNotExist(err) {
		return fmt.Errorf("%s folder not found in %s", actionsFolder, path)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s in %s is not a folder", actionsFolder, path)
	}
	return nil
}

func getCurrentDir() string {
	dir := os.Getenv("NUV_DIR")
	if dir == "" {
		dir, _ = os.Getwd()
		//nolint:errcheck
		os.Setenv("NUV_DIR", dir)
	}
	return dir
}

func printScanUsage() {
	fmt.Println(`Usage:
nuv -scan [options] <dir> <nuv cmd> [args...]

The scan tool visits <dir>/actions and runs the <nuv cmd> on it and all subdirectories recursively.
If the 'actions' folder does not exist, it exits.

Then it invokes the given nuv command with the following arguments:
	- the absolute path of the folder it is processing
	- all the files in the folder that matches the glob (none by default, use '*' for select all)

You can pass a glob pattern with the -g flag to filter the files used as input.

Example:
nuv -scan -g * nuv -js script.js

This results in running the script.js file on the $NUV_DIR/actions folder and all subdirectories. 
For example, if $NUV_DIR/actions contains a subfolder called 'subfolder' with a file called 'bar.js',
the following commands are executed:

	- nuv -js script.js $NUV_DIR/actions
	- nuv -js script.js $NUV_DIR/actions/subfolder bar.js

Options:
  -h		    show help
  -g <pattern>	glob pattern to filter files (default: none, no files are passed to the nuv command)
  -d <dir>      directory to scan (default: $NUV_DIR=PWD)`)
}
