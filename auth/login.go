package auth

import (
	"flag"
	"fmt"

	"github.com/zalando/go-keyring"
)

const usage = `Usage:
nuv login <apihost> [<user>]`

const whiskLoginPath = "/api/v1/web/whisk-system/login"
const defaultUser = "nuvolaris"
const nuvSecretServiceName = "nuvolaris"

func LoginCmd(args []string) error {
	flag.Usage = func() {
		fmt.Println(usage)
	}

	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	fmt.Print("Enter Password: ")
	_, err := AskPassword()
	if err != nil {
		fmt.Println()
		return err
	}
	// url := args[0] + whiskLoginPath
	// user := defaultUser
	// if len(args) > 2 {
	// 	user = args[1]
	// }

	return nil
}

func storeCredentials(creds map[string]string) error {
	for k, v := range creds {
		err := keyring.Set(nuvSecretServiceName, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
