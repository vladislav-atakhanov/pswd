package pswd

import (
	"fmt"
	"os"
	"path/filepath"
	"pswd/internal/crypto"
)

func (p *Pswd) Show(name string, master passwordGetter) (string, error) {
	cipher, err := os.ReadFile(p.Passfile(name))
	if err != nil {
		return "", fmt.Errorf("%s is not in the %s", name, filepath.Base(p.storagePath))
	}

	dir := p.Path(filepath.Dir(name))
	priv, err := p.readPrivateKey(dir)
	if err != nil {
		return "", err
	}
	m, err := master()
	if err != nil {
		return "", err
	}
	plaintext, err := crypto.Decrypt(cipher, priv, m)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func (p *Pswd) readPrivateKey(dir string) ([]byte, error) {
	for {
		pth := p.keysDir(dir)
		if isDir(pth) {
			if c, err := os.ReadFile(p.privateKey(dir)); err == nil {
				return c, nil
			}
		}
		dir = filepath.Dir(dir)
		if dir == filepath.Dir(p.storagePath) {
			return nil, fmt.Errorf("keys not found")
		}
	}
}
