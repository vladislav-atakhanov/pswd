package keys

import (
	"fmt"
	"os"
	"path"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

type KeyStore struct {
	path string
}

func NewKeyStore(storePath string) (*KeyStore, error) {
	if storePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		storePath = path.Join(home, ".keys")
	}
	return &KeyStore{path: storePath}, nil
}

func (ks *KeyStore) keyPath(names ...string) string {
	return path.Join(append([]string{ks.path}, names...)...)
}

func Generate(password string) (priv []byte, pub []byte, err error) {
	return crypto.GenerateKeys(password)
}

func (ks *KeyStore) Exists(id string) bool {
	s, err := os.Stat(ks.keyPath(id, "private.asc"))
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func (ks *KeyStore) Save(id string, priv, pub []byte) (string, error) {
	dir := ks.keyPath(id)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
	}
	if err := os.WriteFile(ks.keyPath(id, "private.asc"), priv, 0644); err != nil {
		return "", fmt.Errorf("save private key: %w", err)
	}
	if err := os.WriteFile(ks.keyPath(id, "public.asc"), pub, 0644); err != nil {
		return "", fmt.Errorf("save public key: %w", err)
	}
	return dir, nil
}

func (ks *KeyStore) Encrypt(id string, plaintext []byte) ([]byte, error) {
	pub, err := os.ReadFile(ks.keyPath(id, "public.asc"))
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(plaintext, pub)
}

func (ks *KeyStore) Decrypt(id string, password string, encData []byte) ([]byte, error) {
	return ks.DecryptLazy(id, func() (string, error) { return password, nil }, encData)
}

func (ks *KeyStore) DecryptLazy(id string, password func() (string, error), encData []byte) ([]byte, error) {
	priv, err := os.ReadFile(ks.keyPath(id, "private.asc"))
	if err != nil {
		return nil, err
	}
	p, err := password()
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(encData, priv, p)
}
