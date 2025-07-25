package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
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
	salt := make([]byte, 16)
	rand.Read(salt)

	aesKey, err := makeAesKey(sharedSecret, salt)
	if err != nil {
		return nil, fmt.Errorf("make aes key: %w", err)
	}

	ciphertext, nonce, err := encryptSync(plaintext, aesKey)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	buf.Write(ephPub)
	buf.Write(salt)
	buf.Write(nonce)
	buf.Write(ciphertext)
	return buf.Bytes(), nil
}

func Decrypt(encData []byte, privateKey []byte, password string) ([]byte, error) {
	const (
		ephPubSize = 32
		saltSize   = 16
		nonceSize  = 12
	)
	if len(encData) < ephPubSize+saltSize+nonceSize {
		return nil, fmt.Errorf("invalid data")
	}

	ephPub := encData[:ephPubSize]
	salt := encData[ephPubSize : ephPubSize+saltSize]
	nonce := encData[ephPubSize+saltSize : ephPubSize+saltSize+nonceSize]
	ciphertext := encData[ephPubSize+saltSize+nonceSize:]

	priv, err := decryptPrivateKey(privateKey, password)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := curve25519.X25519(priv, ephPub)
	if err != nil {
		return nil, fmt.Errorf("invalid ephemeral public key")
	}
	aesKey, err := makeAesKey(sharedSecret, salt)
	if err != nil {
		return nil, fmt.Errorf("make aes key: %w", err)
	}
	return decryptSync(ciphertext, nonce, aesKey)
}

func makeAesKey(sharedSecret []byte, salt []byte) ([]byte, error) {
	h := hkdf.New(sha256.New, sharedSecret, salt, []byte("ecc-encryption"))
	aesKey := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(h, aesKey); err != nil {
		return nil, err
	}
	return aesKey, nil
}
