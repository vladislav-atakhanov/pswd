package pswd

import (
	"os"
	"path"
)

type Pswd struct {
	storagePath string
}

func getStorageDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".pswd"), nil
}

func NewPswd(storagePath string) (*Pswd, error) {
	if storagePath == "" {
		s, err := getStorageDir()
		if err != nil {
			return nil, err
		}
		storagePath = s
	}
	return &Pswd{
		storagePath: storagePath,
	}, nil
}
