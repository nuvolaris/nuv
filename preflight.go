package main

import (
	"errors"
	"os/exec"
)

type preflight struct {
	err error
}

type checkFn func(pd *preflight)

func (p *preflight) check(f checkFn) {

	if p.err != nil {
		// skip if previous check failed
		return
	}
	f(p)
}

// preflightChecks performs preflight checks:
// - curl is installed
// - ssh is installed
func preflightChecks() error {
	trace("preflight checks")
	preflight := preflight{}

	preflight.check(curlInstalled)
	preflight.check(sshInstalled)

	return preflight.err
}

func curlInstalled(p *preflight) {
	trace("check curl is installed")
	_, err := exec.LookPath("curl")
	if err != nil {
		p.err = errors.New("curl is required but is not installed")
	}
}

func sshInstalled(p *preflight) {
	trace("check ssh is installed")
	_, err := exec.LookPath("ssh")
	if err != nil {
		p.err = errors.New("ssh is required but is not installed")
	}
}
