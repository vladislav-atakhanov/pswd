package pswd

import (
	"os"
	"path"

	"github.com/vladislav-atakhanov/pswd/pkg/keys"
)

type Pswd struct {
	storagePath string
	keyStore    *keys.KeyStore
}

func getStorageDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".pswd"), nil
}

func NewPswd(storagePath string, keyStorePath string) (*Pswd, error) {
	if storagePath == "" {
		s, err := getStorageDir()
		if err != nil {
			return nil, err
		}
		storagePath = s
	}

	ks, err := keys.NewKeyStore(keyStorePath)
	if err != nil {
		return nil, err
	}

	return &Pswd{
		storagePath: storagePath,
		keyStore:    ks,
	}, nil
}
