package tools

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type deployCtx struct {
	path               string
	dryRun             bool
	packageCmdExecuted map[string]bool
}

var nuvPackageUpdate = "nuv package update"

func DeployTool() error {
	flag := flag.NewFlagSet("deploy", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Println(`Command to deploy Nuvolaris projects. It takes the path to the project to deploy (must contain 'packages' folder) as argument.

Usage:
	nuv -deploy [Options] <path>

Options:
	-s, --single <string>     Deploy a single action with the given path, either a single file or a directory.
	-w, --watch     Watch for changes and deploy automatically.
	-d, --dry-run   Do not deploy, just print the deployment plan.
	--no-web        Do not upload the web folder to Nuvolaris (if present)`)
	}

	// var watchFlag bool
	// var noWebFlag int
	// var forceFlag bool
	var singleFlag string
	var helpFlag bool
	var dryRunFlag bool
	flag.StringVar(&singleFlag, "s", "", "Deploy a single action with the given path, either a single file or a directory.")
	flag.StringVar(&singleFlag, "single", "", "Deploy a single action with the given path, either a single file or a directory.")
	flag.BoolVar(&helpFlag, "h", false, "Show this help message")
	flag.BoolVar(&helpFlag, "help", false, "Show this help message")
	flag.BoolVar(&dryRunFlag, "d", false, "Do not deploy, just print the deployment plan.")
	flag.BoolVar(&dryRunFlag, "dry-run", false, "Do not deploy, just print the deployment plan.")

	// Parse command line flags
	err := flag.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	if (flag.NArg() == 0) || helpFlag {
		flag.Usage()
		return nil
	}

	rootPath := flag.Arg(0)

	// if there is not "packages" folder from where deploy was called, abort
	if !exists(filepath.Join(rootPath, "packages")) {
		return fmt.Errorf("no 'packages' folder found in the given directory (%s)", rootPath)
	}

	ctx := deployCtx{
		path:               rootPath,
		packageCmdExecuted: make(map[string]bool),
		dryRun:             dryRunFlag,
	}

	if singleFlag != "" {
		action := singleFlag
		if !strings.HasPrefix(action, "packages") {
			action = filepath.Join("packages", action)
			if !exists(filepath.Join(ctx.path, action)) {
				return fmt.Errorf("action %s not found: must be either a file or a directory under packages", action)
			}
		}
		return deploy(ctx, action)
	}
	// walk and deploy
	//scan()

	// watch if requested
	// if  args.watch:
	//     print(">>> Watching:")
	//     watch()

	return nil
}

var supportedMains = []string{"__main__.py", "index.js", "main.js", "main.go"}

func deploy(ctx deployCtx, actionPath string) error {
	fullPath := filepath.Join(ctx.path, actionPath)
	if !exists(fullPath) {
		return fmt.Errorf("action %s not found: must be either a file or a directory", actionPath)
	}
	log.Println("***", filepath.Base(actionPath))

	action, err := checkActionDir(ctx.path, actionPath)
	if err != nil {
		return err
	}

	sp := splitPath(action)
	if len(sp) > 3 {
		action, err = buildAction(ctx, sp[1], sp[2])
		if err != nil {
			return err
		}
	}

	return deployAction(ctx, action)
}

func checkActionDir(rootPath string, actionPath string) (string, error) {
	fullPath := filepath.Join(rootPath, actionPath)
	isActionDir := false
	isActionDirSupported := false
	if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.Mode().IsDir() {
		isActionDir = true
		for _, start := range supportedMains {
			sub := filepath.Join(actionPath, start)
			if exists(filepath.Join(rootPath, sub)) {
				actionPath = sub
				isActionDirSupported = true
				break
			}
		}
	}
	if isActionDir && !isActionDirSupported {
		return "", fmt.Errorf("action %s is a directory but does not contain a supported main file", actionPath)
	}
	return actionPath, nil
}

func deployAction(ctx deployCtx, artifact string) error {
	sp := splitPath(artifact)
	nameType := strings.Split(sp[len(sp)-1], ".")
	name := nameType[0]
	typ := nameType[1]
	packageName := filepath.Base(filepath.Dir(artifact))

	if packageName != "packages" {
		deployPackage(ctx, packageName)
	}

	var toInspect []string
	if typ == "zip" {
		base := filepath.Join(ctx.path, artifact[:len(artifact)-4])

		// TODO: add support for other languages
		toInspect = []string{filepath.Join(base, "/__main__.py"), filepath.Join(base, "/index.js")}
	} else {
		toInspect = []string{artifact}
	}

	args := strings.Join(extractArgs(toInspect), " ")
	action := packageName + "/" + name // the action name, it's not a file path
	fullPath := filepath.Join(ctx.path, artifact)
	if packageName == "packages" {
		action = name
	}
	if !ctx.dryRun {
		cmd := []string{"action", "update", action, fullPath, args}
		err := exec.Command("nuv", cmd...).Run()
		if err != nil {
			log.Println("Error deploying action", name, err)
		}
	} else {
		log.Println("Would run:", "nuv action update", action, fullPath, args)
	}

	return nil
}

func deployPackage(ctx deployCtx, pkg string) {
	// package args
	ppath := filepath.Join(ctx.path, "packages", pkg+".args")
	pargs := strings.Join(extractArgs([]string{ppath}), " ")
	cmd := fmt.Sprintf("%s %s %s", nuvPackageUpdate, pkg, pargs)
	if _, ok := ctx.packageCmdExecuted[cmd]; !ok {
		if !ctx.dryRun {
			err := exec.Command(cmd).Run()
			if err != nil {
				log.Println("Error deploying package", pkg, err)
			}
		} else {
			log.Println("Would run:", cmd)
		}

		ctx.packageCmdExecuted[cmd] = true
	}
}

func extractArgs(files []string) []string {
	res := []string{}
	for _, file := range files {
		if exists(file) {
			f, err := os.Open(file)
			if err != nil {
				log.Println("Error opening file", file, err)
				continue
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#-") {
					res = append(res, strings.TrimSpace(line[2:]))
				}
				if strings.HasPrefix(line, "//-") {
					res = append(res, strings.TrimSpace(line[3:]))
				}
			}

			if err := scanner.Err(); err != nil {
				log.Println("Error reading file", file, err)
			}
		}
	}
	return res
}

func buildAction(ctx deployCtx, pkg string, action string) (string, error) {
	if !ctx.dryRun {
		err := exec.Command("nuv", "ide", "util", "action", fmt.Sprintf("A=%s/%s", pkg, action)).Run()
		if err != nil {
			return "", fmt.Errorf("error building action %s/%s: %v", pkg, action, err)
		}
	} else {
		log.Println("Would run: nuv ide util action A=" + pkg + "/" + action)
	}
	return fmt.Sprintf("packages/%s/%s.zip", pkg, action), nil
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func splitPath(path string) []string {
	dir, last := filepath.Split(path)
	if dir == "" {
		return []string{last}
	}
	return append(splitPath(filepath.Clean(dir)), last)
}
