package pswd

import (
	"os"
	"path/filepath"

	"github.com/vladislav-atakhanov/pswd/pkg/keys"
)

func (p *Pswd) Insert(name string, password string) (string, error) {
	id, err := p.getKeyId(name)
	if err != nil {
		return "", err
	}
	return p.InsertWithKey(id, name, password)
}

func (p *Pswd) InsertWithKey(id string, name string, password string) (string, error) {
	dir := p.Path(filepath.Dir(name))
	cipher, err := keys.Encrypt(id, []byte(password))
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
