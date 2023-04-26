package tools

import (
	"flag"
	"fmt"
	"log"
	osuser "os/user"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

var (
	help bool

	user       string
	addr       string
	port       uint
	key        string
	cmd        string
	pass       bool
	passphrase bool
	timeout    time.Duration
	agent      bool

	err error
)

func SshTool() error {
	flag.Usage = printSSHUsage

	flag.BoolVar(&help, "help", false, "print this help message.")
	flag.BoolVar(&help, "h", false, "print this help message.")

	flag.StringVar(&user, "user", "root", "ssh user.")
	flag.StringVar(&addr, "ip", "127.0.0.1", "machine ip address.")
	flag.UintVar(&port, "port", 22, "ssh port number.")
	flag.StringVar(&key, "key", "", "private key path.")
	flag.StringVar(&cmd, "cmd", "", "command to run.")
	flag.BoolVar(&pass, "pass", false, "ask for ssh password instead of private key.")
	flag.BoolVar(&agent, "agent", false, "use ssh agent for authentication (unix systems only).")
	flag.BoolVar(&passphrase, "passphrase", false, "ask for private key passphrase.")
	flag.DurationVar(&timeout, "timeout", 0, "interrupt a command with SIGINT after a given timeout (0 means no timeout)")

	// Parse command line flags
	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	// if key is not provided, use default ~/.ssh/id_rsa
	if !isFlagPassed("key") {
		k, err := defaultSshKeyPath()
		if err != nil {
			return err
		}
		key = k
	}

	// if user is not provided, use current user
	if !isFlagPassed("user") {
		user = defaultSshUser()
	}

	return nil
}

func printSSHUsage() {
	fmt.Print(`Usage:
nuv -ssh [options]

-h, --help 		   print this help message.
-u, --user 		   the ssh user (default: current user)
`)
}

func defaultSshUser() string {
	usr, err := osuser.Current()
	if err != nil {
		log.Println("couldn't determine current user. Defaulting to 'root'")
		return "root"
	}
	return usr.Username
}

func defaultSshKeyPath() (string, error) {
	homessh, err := homedir.Expand("~/.ssh")
	if err != nil {
		return "", err
	}
	return filepath.Join(homessh, "id_rsa"), nil
}
