package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

func preparePassword(password string) []byte {
	key := sha256.Sum256([]byte(password))
	aesKey := key[:]
	return aesKey
}

func encryptPrivateKey(privateKey []byte, password string) ([]byte, error) {
	private, nonce, err := encryptSync(privateKey, preparePassword(password))
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	buf.Write(nonce)
	buf.Write(private)
	return buf.Bytes(), nil
}

func decryptPrivateKey(encData []byte, password string) ([]byte, error) {
	if len(encData) < 12 {
		return nil, fmt.Errorf("encrypted data too short")
	}
	nonce := encData[:12]
	ciphertext := encData[12:]
	return decryptSync(ciphertext, nonce, preparePassword(password))
}
