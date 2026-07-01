package crypto

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrWrongPassword  = errors.New("wrong password or corrupted key file")
	ErrInvalidKeyFile = errors.New("invalid encrypted key length")
)

const (
	SaltLen  = 16
	NonceLen = chacha20poly1305.NonceSizeX
	KeyLen   = 32
)

func EncryptPrivateKey(priv [32]byte, password []byte) ([]byte, error) {
	salt := make([]byte, SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	derivedKey := argon2.IDKey(password, salt, 3, 64*1024, 4, KeyLen)
	mem.Lock(derivedKey)
	defer mem.Unlock(derivedKey)
	defer mem.ZeroBytes(derivedKey)

	aead, err := chacha20poly1305.NewX(derivedKey)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aead.Seal(nil, nonce, priv[:], nil)

	result := make([]byte, 0, SaltLen+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func DecryptPrivateKey(data []byte, password []byte) ([32]byte, error) {
	if len(data) < SaltLen+NonceLen+KeyLen+chacha20poly1305.Overhead {
		return [32]byte{}, ErrInvalidKeyFile
	}

	salt := data[:SaltLen]
	nonce := data[SaltLen : SaltLen+NonceLen]
	ciphertext := data[SaltLen+NonceLen:]

	derivedKey := argon2.IDKey(password, salt, 3, 64*1024, 4, KeyLen)
	mem.Lock(derivedKey)
	defer mem.Unlock(derivedKey)
	defer mem.ZeroBytes(derivedKey)

	aead, err := chacha20poly1305.NewX(derivedKey)
	if err != nil {
		return [32]byte{}, err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return [32]byte{}, ErrWrongPassword
	}

	var priv [32]byte
	copy(priv[:], plaintext)
	return priv, nil
}
