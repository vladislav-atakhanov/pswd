package keys

import (
	"fmt"
	"os"
	"path"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

func Generate(password string) (priv []byte, pub []byte, err error) {
	return crypto.GenerateKeys(password)
}

func Save(id string, priv, pub []byte) error {
	if err := os.MkdirAll(keyPath(id), 0700); err != nil {
		return fmt.Errorf("create keys dir: %w", err)
	}
	if err := os.WriteFile(keyPath(id, "private.asc"), priv, 0644); err != nil {
		return fmt.Errorf("save private key: %w", err)
	}
	if err := os.WriteFile(keyPath(id, "public.asc"), pub, 0644); err != nil {
		return fmt.Errorf("save public key: %w", err)
	}
	return nil
}

func keyPath(names ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(append([]string{home, ".keys"}, names...)...)
}

func Encrypt(id string, plaintext []byte) ([]byte, error) {
	pub, err := os.ReadFile(keyPath(id, "public.asc"))
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(plaintext, pub)
}
func Decrypt(id string, password string, encData []byte) ([]byte, error) {
	return DecryptLazy(id, func() (string, error) { return password, nil }, encData)
}
func DecryptLazy(id string, password func() (string, error), encData []byte) ([]byte, error) {
	priv, err := os.ReadFile(keyPath(id, "private.asc"))
	if err != nil {
		return nil, err
	}
	p, err := password()
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(encData, priv, p)
}
