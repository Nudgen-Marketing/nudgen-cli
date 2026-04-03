package config

import (
	"log"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "nudgen-cli"
	userAccount = "pat-token" // We only need one token globally right now
)

// SaveToken saves the Personal Access Token securely in the OS's native keychain.
func SaveToken(token string) error {
	return keyring.Set(serviceName, userAccount, token)
}

// GetToken retrieves the stored PAT.
func GetToken() (string, error) {
	return keyring.Get(serviceName, userAccount)
}

// ClearToken removes the stored PAT (logout).
func ClearToken() error {
	return keyring.Delete(serviceName, userAccount)
}

// CheckToken is a helper to see if a token exists
func IsLoggedIn() bool {
	_, err := GetToken()
	return err == nil
}

// EnsureLoggedIn enforces that the user must log in before running a command.
func EnsureLoggedIn() {
	if !IsLoggedIn() {
		log.Fatal("You are not logged in. Please run `nudgen login` first.")
	}
}
