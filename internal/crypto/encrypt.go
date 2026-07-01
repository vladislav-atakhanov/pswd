package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"golang.org/x/crypto/chacha20poly1305"
)

type countingWriter struct {
	w     io.Writer
	count int
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.count += n
	return n, err
}

func EncryptStream(out io.Writer, plainText io.Reader, publicKeys [][32]byte) (int, error) {
	cw := &countingWriter{w: out}

	plain, err := io.ReadAll(plainText)
	if err != nil {
		return 0, err
	}
	defer mem.ZeroBytes(plain)

	dataKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
		return 0, err
	}
	mem.Lock(dataKey)
	defer mem.Unlock(dataKey)
	defer mem.ZeroBytes(dataKey)

	ephemeralPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return 0, err
	}
	ephemeralPubBytes := ephemeralPriv.PublicKey().Bytes()

	if _, err := cw.Write(ephemeralPubBytes); err != nil {
		return cw.count, err
	}

	numDevices := uint16(len(publicKeys))
	if err := binary.Write(cw, binary.BigEndian, numDevices); err != nil {
		return cw.count, err
	}

	for _, pubBytes := range publicKeys {
		devicePub, err := ecdh.X25519().NewPublicKey(pubBytes[:])
		if err != nil {
			return cw.count, err
		}

		sharedSecret, err := ephemeralPriv.ECDH(devicePub)
		if err != nil {
			return cw.count, err
		}

		aead, err := chacha20poly1305.New(sharedSecret)
		if err != nil {
			return cw.count, err
		}

		nonce := make([]byte, aead.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return cw.count, err
		}

		encryptedKey := aead.Seal(nil, nonce, dataKey, nil)

		if _, err := cw.Write(nonce); err != nil {
			return cw.count, err
		}
		if _, err := cw.Write(encryptedKey); err != nil {
			return cw.count, err
		}
	}

	aead, err := chacha20poly1305.New(dataKey)
	if err != nil {
		return cw.count, err
	}

	mainNonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, mainNonce); err != nil {
		return cw.count, err
	}

	ciphertext := aead.Seal(nil, mainNonce, plain, nil)

	if _, err := cw.Write(mainNonce); err != nil {
		return cw.count, err
	}
	if _, err := cw.Write(ciphertext); err != nil {
		return cw.count, err
	}

	return cw.count, nil
}
