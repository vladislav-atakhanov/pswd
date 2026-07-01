package crypto

import (
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/chacha20"
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

	dataKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
		return 0, err
	}

	ephemeralPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return 0, err
	}
	ephemeralPubBytes := ephemeralPriv.PublicKey().Bytes()

	if _, err := cw.Write(ephemeralPubBytes); err != nil {
		return cw.count, err
	}

	macKey := sha256.Sum256(append(dataKey, []byte("mac")...))
	hmacSigner := hmac.New(sha256.New, macKey[:])

	teeOut := io.MultiWriter(cw, hmacSigner)

	for _, pubBytes := range publicKeys {
		devicePub, err := ecdh.X25519().NewPublicKey(pubBytes[:])
		if err != nil {
			return cw.count, err
		}

		sharedSecret, err := ephemeralPriv.ECDH(devicePub)
		if err != nil {
			return cw.count, err
		}

		keyNonce := make([]byte, chacha20.NonceSize)
		cipherKey, err := chacha20.NewUnauthenticatedCipher(sharedSecret, keyNonce)
		if err != nil {
			return cw.count, err
		}

		encryptedKey := make([]byte, 32)
		cipherKey.XORKeyStream(encryptedKey, dataKey)

		if _, err := teeOut.Write(encryptedKey); err != nil {
			return cw.count, err
		}
	}

	mainNonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(rand.Reader, mainNonce); err != nil {
		return cw.count, err
	}
	if _, err := teeOut.Write(mainNonce); err != nil {
		return cw.count, err
	}

	mainCipher, err := chacha20.NewUnauthenticatedCipher(dataKey, mainNonce)
	if err != nil {
		return cw.count, err
	}

	buf := make([]byte, 4096)
	for {
		n, err := plainText.Read(buf)
		if n > 0 {
			mainCipher.XORKeyStream(buf[:n], buf[:n])
			if _, err := teeOut.Write(buf[:n]); err != nil {
				return cw.count, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return cw.count, err
		}
	}

	signature := hmacSigner.Sum(nil)
	if _, err := cw.Write(signature); err != nil {
		return cw.count, err
	}

	return cw.count, nil
}
