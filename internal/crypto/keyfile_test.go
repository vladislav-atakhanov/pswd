package crypto

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/chacha20poly1305"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	password := []byte("correct horse battery staple")

	encrypted, err := EncryptPrivateKey(key, password)
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	decrypted, err := DecryptPrivateKey(encrypted, password)
	if err != nil {
		t.Fatalf("DecryptPrivateKey: %v", err)
	}

	if !bytes.Equal(key[:], decrypted[:]) {
		t.Fatal("round trip mismatch")
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	encrypted, err := EncryptPrivateKey(key, []byte("correct password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	_, err = DecryptPrivateKey(encrypted, []byte("wrong password"))
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestDecryptCorruptedCiphertext(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	encrypted, err := EncryptPrivateKey(key, []byte("password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	encrypted[len(encrypted)-1] ^= 0xff

	_, err = DecryptPrivateKey(encrypted, []byte("password"))
	if err == nil {
		t.Fatal("expected error for corrupted ciphertext, got nil")
	}
}

func TestDecryptTruncatedData(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	encrypted, err := EncryptPrivateKey(key, []byte("password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	_, err = DecryptPrivateKey(encrypted[:10], []byte("password"))
	if err == nil {
		t.Fatal("expected error for truncated data, got nil")
	}
}

func TestEmptyPassword(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	encrypted, err := EncryptPrivateKey(key, []byte{})
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	decrypted, err := DecryptPrivateKey(encrypted, []byte{})
	if err != nil {
		t.Fatalf("DecryptPrivateKey: %v", err)
	}

	if !bytes.Equal(key[:], decrypted[:]) {
		t.Fatal("round trip mismatch with empty password")
	}
}

func TestLongPassword(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	password := make([]byte, 1000)

	encrypted, err := EncryptPrivateKey(key, password)
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	decrypted, err := DecryptPrivateKey(encrypted, password)
	if err != nil {
		t.Fatalf("DecryptPrivateKey: %v", err)
	}

	if !bytes.Equal(key[:], decrypted[:]) {
		t.Fatal("round trip mismatch with long password")
	}
}

func TestEncryptUniqueOutput(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	enc1, err := EncryptPrivateKey(key, []byte("password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	enc2, err := EncryptPrivateKey(key, []byte("password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	if bytes.Equal(enc1, enc2) {
		t.Fatal("two encryptions with same key and password produced identical output")
	}
}

func TestEncryptOutputLength(t *testing.T) {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	encrypted, err := EncryptPrivateKey(key, []byte("password"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}

	expected := SaltLen + VerifierLen + NonceLen + 32 + chacha20poly1305.Overhead
	if len(encrypted) != expected {
		t.Fatalf("expected length %d, got %d", expected, len(encrypted))
	}
}

func TestDecryptInvalidLength(t *testing.T) {
	_, err := DecryptPrivateKey([]byte{1, 2, 3}, []byte{})
	if err == nil {
		t.Fatal("expected error for invalid length, got nil")
	}
}
