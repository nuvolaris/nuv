package auth

import (
	"errors"
	"syscall"

	"golang.org/x/term"
)

type PasswordReader interface {
	ReadPassword() (string, error)
}

type StdInPasswordReader struct{}

func (r StdInPasswordReader) ReadPassword() (string, error) {
	pwd, error := term.ReadPassword(syscall.Stdin)
	return string(pwd), error
}

var pwdReader PasswordReader = StdInPasswordReader{}

func AskPassword() (string, error) {
	pwd, err := pwdReader.ReadPassword()
	if err != nil {
		return "", err
	}
	if len(pwd) == 0 {
		return "", errors.New("password is empty")
	}
	return pwd, nil
}
