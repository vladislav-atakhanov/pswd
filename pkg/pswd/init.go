package pswd

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

func (p *Pswd) IsInit(d string) bool {
	if isFile(p.privateKey(p.storagePath, d)) && isFile(p.publicKey(p.storagePath, d)) {
		return true
	}
	return false
}
func (p *Pswd) Init(dir string, new, old func() (string, error)) (string, []string, error) {
	if !p.IsInit(dir) {
		d, err := p.init(dir, new)
		if err != nil {
			return "", nil, err
		}
		return d, nil, nil
	}
	return p.reinit(dir, new, old)
}

func (p *Pswd) init(dir string, master func() (string, error)) (string, error) {
	if p.IsInit(dir) {
		return "", fmt.Errorf("keys are exist")
	}
	d := p.Path(dir)
	if err := os.MkdirAll(p.keysDir(d), 0700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
	}
	password, err := master()
	if err != nil {
		return "", err
	}
	priv, pub, err := crypto.GenerateKeys(password)
	if err != nil {
		return "", fmt.Errorf("generate keys: %w", err)
	}
	if err := p.saveKeys(d, priv, pub); err != nil {
		return "", err
	}
	return d, nil
}

func (p *Pswd) reinit(dir string, new, old passwordGetter) (string, []string, error) {
	d := p.Path(dir)

	var priv, pub []byte

	files, err := walk(d, func(fp string, de fs.DirEntry) bool {
		if fp == p.storagePath {
			return true
		}
		if strings.HasPrefix(de.Name(), ".") {
			return false
		}
		if de.IsDir() {
			return true
		}
		return strings.HasSuffix(de.Name(), ".asc")
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	var oldPass string

	names := []string{}
	for _, f := range files {
		name := p.passfileToName(f)
		if oldPass == "" {
			oldPass, err = old()
			if err != nil {
				return "", nil, err
			}
		}
		password, err := p.Show(name, func() (string, error) { return oldPass, nil })
		if err != nil {
			return "", nil, err
		}
		if priv == nil {
			newPass, err := new()
			if err != nil {
				return "", nil, err
			}
			priv, pub, err = crypto.GenerateKeys(newPass)
			if err != nil {
				return "", nil, fmt.Errorf("generate keys: %w", err)
			}
		}
		_, err = p.InsertWithKey(pub, name, password)
		if err != nil {
			return "", nil, err
		}
		names = append(names, name)
	}

	if err := p.saveKeys(d, priv, pub); err != nil {
		return "", nil, err
	}
	return d, names, nil
}

func (p *Pswd) saveKeys(d string, priv, pub []byte) error {
	if err := os.WriteFile(p.privateKey(d), priv, 0644); err != nil {
		return fmt.Errorf("save private key: %w", err)
	}
	if err := os.WriteFile(p.publicKey(d), pub, 0644); err != nil {
		return fmt.Errorf("save public key: %w", err)
	}
	return nil
}
