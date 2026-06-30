package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
)

func PublicKeyFromPrivate(privateBytes [32]byte) ([32]byte, error) {
	var publicBytes [32]byte

	priv, err := ecdh.X25519().NewPrivateKey(privateBytes[:])
	if err != nil {
		return publicBytes, err
	}

	copy(publicBytes[:], priv.PublicKey().Bytes())
	return publicBytes, nil
}

func GenerateKeys() (privateBytes [32]byte, publicBytes [32]byte, err error) {
	privKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return privateBytes, publicBytes, err
	}
	copy(privateBytes[:], privKey.Bytes())
	copy(publicBytes[:], privKey.PublicKey().Bytes())
	return privateBytes, publicBytes, nil
}
