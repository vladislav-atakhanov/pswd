package pswd

import (
	"fmt"
	"os"
	"pswd/internal/crypto"
)

func (p *Pswd) IsInit(d string) bool {
	if isFile(p.privateKey(p.storagePath, d)) && isFile(p.publicKey(p.storagePath, d)) {
		return true
	}
	return false
}

func (p *Pswd) Init(dir string, password string) (string, error) {
	if p.IsInit(dir) {
		return "", fmt.Errorf("keys are exist")
	}
	d := p.Path(dir)
	if err := os.MkdirAll(p.keysDir(d), 0700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
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
func (p *Pswd) saveKeys(d string, priv, pub []byte) error {
	if err := os.WriteFile(p.privateKey(d), priv, 0644); err != nil {
		return fmt.Errorf("save private key: %w", err)
	}
	if err := os.WriteFile(p.publicKey(d), pub, 0644); err != nil {
		return fmt.Errorf("save public key: %w", err)
	}
	return nil
}
