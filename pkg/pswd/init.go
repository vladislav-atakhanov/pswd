package pswd

import (
	"fmt"
	"os"
	"path"
	"pswd/internal/crypto"
)

func (p *Pswd) Init(dir string, password string) (string, error) {
	d := path.Join(p.storagePath, dir)

	if isFile(p.privateKey(d)) && isFile(p.publicKey(d)) {
		return "", fmt.Errorf("keys are exist")
	}
	if err := os.MkdirAll(p.keysDir(d), 0700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
	}
	priv, pub, err := crypto.GenerateKeys(password)
	if err != nil {
		return "", fmt.Errorf("generate keys: %w", err)
	}
	if err := os.WriteFile(p.privateKey(d), priv, 0644); err != nil {
		return "", fmt.Errorf("save private key: %w", err)
	}
	if err := os.WriteFile(p.publicKey(d), pub, 0644); err != nil {
		return "", fmt.Errorf("save public key: %w", err)
	}
	return d, nil
}
