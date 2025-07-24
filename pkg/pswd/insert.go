package pswd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pswd/internal/crypto"
)

func (p *Pswd) Insert(name string, password string) (string, error) {
	dir := path.Join(p.storagePath, filepath.Dir(name))
	pub, err := p.readPublicKey(dir)
	if err != nil {
		return "", err
	}
	cipher, err := crypto.Encrypt([]byte(password), pub)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	passfile := path.Join(p.storagePath, name) + ".asc"
	if err := os.WriteFile(passfile, cipher, 0644); err != nil {
		return "", err
	}
	return passfile, nil
}

func (p *Pswd) readPublicKey(dir string) ([]byte, error) {
	for {
		pth := p.keysDir(dir)
		if isDir(pth) {
			if c, err := os.ReadFile(p.publicKey(dir)); err == nil {
				return c, nil
			}
		}
		dir = filepath.Dir(dir)
		if dir == filepath.Dir(p.storagePath) {
			return nil, fmt.Errorf("keys not found")
		}
	}
}
