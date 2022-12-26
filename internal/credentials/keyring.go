package credentials

import (
	"fmt"

	"github.com/99designs/keyring"
)

type keyringManager struct {
	ring *keyring.Keyring
}

func (k *keyringManager) RequiredEnv() []string {
	return nil
}

func (k *keyringManager) Resolve(path string, key string) (string, error) {
	kr, err := keyring.Open(keyring.Config{
		ServiceName: path,
	})
	if err != nil {
		return "", fmt.Errorf("failed to open keyring: %w", err)
	}
	keys, err := kr.Keys()
	if err != nil {
		return "", fmt.Errorf("faild to list keys: %w", err)
	}
	fmt.Println(keys)
	secret, err := kr.Get(key)
	if err != nil {
		return "", fmt.Errorf("failed to read key %q from keyring: %w", key, err)
	}
	return string(secret.Data), nil
}

func (k *keyringManager) Store(path string, key string, value string) error {
	kr, err := keyring.Open(keyring.Config{
		ServiceName: path,
	})
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}
	secret := keyring.Item{
		Key:  key,
		Data: []byte(value),
	}
	err = kr.Set(secret)
	if err != nil {
		return fmt.Errorf("failed to set key %q in keyring: %w", key, err)
	}
	return nil
}
