package main

import (
	"encoding/base64"
	"encoding/binary"
	"testing"
)

func makeToken(pub []byte, name string) string {
	buf := make([]byte, 0, 32+2+len(name))
	buf = append(buf, pub...)
	var nameLen [2]byte
	binary.BigEndian.PutUint16(nameLen[:], uint16(len(name)))
	buf = append(buf, nameLen[:]...)
	buf = append(buf, name...)
	return base64.URLEncoding.EncodeToString(buf)
}

func TestParsePublicKey(t *testing.T) {
	var pub [32]byte
	for i := range pub {
		pub[i] = byte(i)
	}
	name := "test-device"

	token := makeToken(pub[:], name)
	gotName, gotPub, err := parsePublicKey(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != name {
		t.Fatalf("expected name %q, got %q", name, gotName)
	}
	if gotPub != pub {
		t.Fatalf("expected pub %x, got %x", pub, gotPub)
	}
}

func TestParsePublicKeyInvalidBase64(t *testing.T) {
	_, _, err := parsePublicKey("!!!invalid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestParsePublicKeyTooShort(t *testing.T) {
	token := base64.URLEncoding.EncodeToString([]byte{1, 2, 3})
	_, _, err := parsePublicKey(token)
	if err == nil {
		t.Fatal("expected error for short token, got nil")
	}
}

func TestParsePublicKeyNameLenExceedsToken(t *testing.T) {
	var pub [32]byte
	buf := make([]byte, 32+2)
	copy(buf, pub[:])
	binary.BigEndian.PutUint16(buf[32:34], 65535)

	token := base64.URLEncoding.EncodeToString(buf)
	_, _, err := parsePublicKey(token)
	if err == nil {
		t.Fatal("expected error for nameLen exceeding token, got nil")
	}
}

func TestParsePublicKeyNoName(t *testing.T) {
	var pub [32]byte
	// 34 bytes: pub(32) + nameLen(2) with value 0
	buf := make([]byte, 34)
	copy(buf, pub[:])
	binary.BigEndian.PutUint16(buf[32:34], 0)
	token := base64.URLEncoding.EncodeToString(buf)

	gotName, _, err := parsePublicKey(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "" {
		t.Fatalf("expected empty name, got %q", gotName)
	}
}
