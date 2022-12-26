package credentials

import (
	"fmt"
	"os"
	"strings"
)

type CredentialManager interface {
	RequiredEnv() []string
	Resolve(path string, key string) (string, error)
	Store(path string, key string, value string) error
}

var credentialManagers = map[string]CredentialManager{
	"vault":   &vault{},
	"keyring": &keyringManager{},
}

type Secret struct {
	Manager string
	Path    string
	Key     string
}

func parseSecret(rawSecret string) (Secret, error) {
	secret := Secret{}
	keys := strings.Split(rawSecret, ":")
	if len(keys) != 3 {
		return Secret{}, fmt.Errorf("key has wrong format, should be 3 separated strings, found %d", len(keys))
	}
	secret.Manager = keys[0]
	secret.Path = keys[1]
	secret.Key = keys[2]
	return secret, nil
}

// Resolve resolves a secret with format manager:path:key
// Example: vault:my_kv/secret:key
func Resolve(secretPath string) (string, error) {
	secret, err := parseSecret(secretPath)
	if err != nil {
		return "", err
	}
	credentialManager, exists := credentialManagers[secret.Manager]
	if !exists {
		return "", fmt.Errorf("secret manager %q not supported", secret.Manager)
	}
	for _, env := range credentialManager.RequiredEnv() {
		if os.Getenv(env) == "" {
			return "", fmt.Errorf("required env variable %q for secret manager %q is missing", env, secret.Manager)
		}
	}
	return credentialManager.Resolve(secret.Path, secret.Key)
}

func Store(secretPath string, data string) error {
	secret, err := parseSecret(secretPath)
	if err != nil {
		return err
	}
	credentialManager, exists := credentialManagers[secret.Manager]
	if !exists {
		return fmt.Errorf("secret manager %q not supported", secret.Manager)
	}
	for _, env := range credentialManager.RequiredEnv() {
		if os.Getenv(env) == "" {
			return fmt.Errorf("required env variable %q for secret manager %q is missing", env, secret.Manager)
		}
	}
	return credentialManager.Store(secret.Path, secret.Key, data)
}
