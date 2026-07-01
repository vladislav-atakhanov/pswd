package main

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

func TestRemoveDevice(t *testing.T) {
	priv1, pub1, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	priv2, pub2, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with device1
	v := vault.New()
	if err := v.InitDevice(pub1, "device1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("password1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("password2"), "entry2"); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Open with priv1, add device2
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := vault.Open(mf, mf.Len(), priv1)
	if err != nil {
		t.Fatalf("open with priv1: %v", err)
	}
	if err := v2.AddDevice(pub2, "device2", mf, priv1); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Open with priv2, remove device1
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := vault.Open(mf, mf.Len(), priv2)
	if err != nil {
		t.Fatalf("open with priv2: %v", err)
	}
	if err := v3.RemoveDevice(pub1, mf, priv2); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v3.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Open with priv2 — must work
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v4, err := vault.Open(mf, mf.Len(), priv2)
	if err != nil {
		t.Fatalf("open with priv2 after removal: %v", err)
	}
	if len(v4.Content) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(v4.Content))
	}

	// Open with priv1 — must fail (device removed)
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	_, err = vault.Open(mf, mf.Len(), priv1)
	if err == nil {
		t.Fatal("expected error opening with removed device, got nil")
	}
}
