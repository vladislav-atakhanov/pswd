package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func GenerateKeys(password string) (priv []byte, pub []byte, err error) {
	var privKey [32]byte
	if _, err := rand.Read(privKey[:]); err != nil {
		return nil, nil, err
	}
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privKey)
	var pb []byte = pubKey[:]
	privateKey, err := encryptPrivateKey(privKey[:], password)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, pb, nil
}

func Encrypt(plaintext []byte, pub []byte) ([]byte, error) {
	ephPriv := make([]byte, 32)
	if _, err := rand.Read(ephPriv); err != nil {
		return nil, err
	}

	ephPub, err := curve25519.X25519(ephPriv, curve25519.Basepoint)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := curve25519.X25519(ephPriv, pub[:])
	if err != nil {
		return nil, fmt.Errorf("invalid recipient public key")
	}
	aesKey := sharedSecret[:32]
	ciphertext, nonce, err := encryptSync(plaintext, aesKey)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	buf.Write(ephPub)
	buf.Write(nonce)
	buf.Write(ciphertext)
	return buf.Bytes(), nil
}

func Decrypt(encData []byte, privateKey []byte, password string) ([]byte, error) {
	if len(encData) < 32+12 {
		return nil, fmt.Errorf("encrypted data too short")
	}

	ephPub := encData[:32]
	nonce := encData[32 : 32+12]
	ciphertext := encData[32+12:]

	priv, err := decryptPrivateKey(privateKey, password)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := curve25519.X25519(priv, ephPub)
	if err != nil {
		return nil, fmt.Errorf("invalid ephemeral public key")
	}
	aesKey := sharedSecret[:32]
	return decryptSync(ciphertext, nonce, aesKey)
}
