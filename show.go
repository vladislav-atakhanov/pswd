package pswd

import (
	"fmt"
	"os"
	"path/filepath"
)

func (p *Pswd) Show(name string, master string) (string, error) {
	return p.ShowLazy(name, func(_ string) (string, error) {
		return master, nil
	})
}

func (p *Pswd) ShowLazy(name string, master func(key string) (string, error)) (string, error) {
	cipher, err := os.ReadFile(p.Passfile(name))
	if err != nil {
		return "", fmt.Errorf("%s is not in the %s", name, filepath.Base(p.storagePath))
	}

	id, err := p.getKeyId(name)
	if err != nil {
		return "", err
	}

	plaintext, err := p.keyStore.DecryptLazy(id, func() (string, error) {
		return master(id)
	}, cipher)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
