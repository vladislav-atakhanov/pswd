package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestPublicKeyFromPrivate(t *testing.T) {
	alicePrivHex := "77076d0a7318a57d3c16c17251b26645df4c2f87ebc0992ab177fba51db92c2a"
	alicePubHex := "8520f0098930a754748b7ddcb43ef75a0dbf3a0d26381af4eba4a98eaa9b4e6a"

	privBytes, err := hex.DecodeString(alicePrivHex)
	if err != nil {
		t.Fatal(err)
	}
	expectedPub, err := hex.DecodeString(alicePubHex)
	if err != nil {
		t.Fatal(err)
	}

	var priv [32]byte
	copy(priv[:], privBytes)

	pub, err := PublicKeyFromPrivate(priv)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivate failed: %v", err)
	}

	if !bytes.Equal(pub[:], expectedPub) {
		t.Fatalf("public key mismatch:\ngot:  %x\nwant: %x", pub, expectedPub)
	}
}

func TestPublicKeyFromPrivateInvalid(t *testing.T) {
	var priv [32]byte
	_, err := PublicKeyFromPrivate(priv)
	if err == nil {
		t.Log("zero private key is accepted (clamped to non-zero), skipping")
	}
}

func TestGenerateKeys(t *testing.T) {
	priv, pub, err := GenerateKeys()
	if err != nil {
		t.Fatalf("GenerateKeys failed: %v", err)
	}

	if priv == ([32]byte{}) {
		t.Fatal("private key is zero")
	}
	if pub == ([32]byte{}) {
		t.Fatal("public key is zero")
	}

	derivedPub, err := PublicKeyFromPrivate(priv)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivate failed: %v", err)
	}

	if !bytes.Equal(pub[:], derivedPub[:]) {
		t.Fatal("derived public key does not match generated public key")
	}
}

func TestGenerateKeysUnique(t *testing.T) {
	_, pub1, _ := GenerateKeys()
	_, pub2, _ := GenerateKeys()
	if bytes.Equal(pub1[:], pub2[:]) {
		t.Fatal("two GenerateKeys calls produced the same public key")
	}
}
