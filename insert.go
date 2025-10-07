package pswd

import (
	"fmt"
	"os"
	"path/filepath"
)

func (p *Pswd) Insert(name string, password string) (string, error) {
	if _, err := p.Type(name); err == nil {
		return "", fmt.Errorf("already exists")
	}
	id, err := p.getKeyId(name)
	if err != nil {
		return "", err
	}
	return p.InsertWithKey(id, name, password)
}

func (p *Pswd) InsertWithKey(id string, name string, password string) (string, error) {
	dir := p.Path(filepath.Dir(name))
	cipher, err := p.keyStore.Encrypt(id, []byte(password))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	passfile := p.Passfile(name)
	if err := os.WriteFile(passfile, cipher, 0644); err != nil {
		return "", err
	}
	return passfile, nil
}
