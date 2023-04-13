package auth

import (
	"errors"
	"testing"
)

type stubPasswordReader struct {
	Password    string
	ReturnError bool
}

func (pr stubPasswordReader) ReadPassword() (string, error) {
	if pr.ReturnError {
		return "", errors.New("stubbed error")
	}
	return pr.Password, nil
}

func TestAskPassword(t *testing.T) {
	t.Run("error: returns error when password reader returns error", func(t *testing.T) {
		oldPwdReader := pwdReader
		pwdReader = stubPasswordReader{ReturnError: true}

		result, err := AskPassword()

		pwdReader = oldPwdReader
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "stubbed error" {
			t.Errorf("Expected error to be 'stubbed error', got %s", err.Error())
		}

		if result != "" {
			t.Errorf("Expected empty string, got %s", result)
		}
	})

	t.Run("success: returns password correctly", func(t *testing.T) {
		oldPwdReader := pwdReader
		pwdReader = stubPasswordReader{Password: "a password", ReturnError: false}

		result, err := AskPassword()

		pwdReader = oldPwdReader

		if err != nil {
			t.Errorf("Expected no error, got %s", err.Error())
		}
		if result != "a password" {
			t.Errorf("Expected 'a password', got %s", result)
		}
	})
}
