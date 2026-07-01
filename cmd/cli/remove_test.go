package main

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

func TestRemoveOneOfTwo(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// First session: create vault with 2 entries
	v := vault.New()
	if err := v.InitDevice(pub, "test-device"); err != nil {
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

	// Second session: open, remove one, save
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	if len(v2.Content) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(v2.Content))
	}

	var id uuid.V4
	for id = range v2.Content {
		break
	}
	if err := v2.Remove(id); err != nil {
		t.Fatal(err)
	}

	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Third session: reopen and verify
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("second open (after remove+save): %v", err)
	}
	if len(v3.Content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.Content))
	}
}

func TestRemoveAll(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with 2 entries
	v := vault.New()
	if err := v.InitDevice(pub, "test-device"); err != nil {
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

	// Open and remove both
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	for id := range v2.Content {
		if err := v2.Remove(id); err != nil {
			t.Fatal(err)
		}
	}
	if len(v2.Content) != 0 {
		t.Fatalf("expected 0 entries after remove all, got %d", len(v2.Content))
	}

	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Reopen — must work (empty vault)
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after remove all: %v", err)
	}
	if len(v3.Content) != 0 {
		t.Fatalf("expected empty vault, got %d", len(v3.Content))
	}
}

func TestRemoveThenAdd(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with 1 entry
	v := vault.New()
	if err := v.InitDevice(pub, "test-device"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("password1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Open, remove the only entry, add a new one
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	for id := range v2.Content {
		if err := v2.Remove(id); err != nil {
			t.Fatal(err)
		}
		break
	}
	if err := v2.Add(strings.NewReader("newpassword"), "newentry"); err != nil {
		t.Fatal(err)
	}

	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Reopen — must have 1 entry with new content
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := vault.Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after remove+add: %v", err)
	}
	if len(v3.Content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.Content))
	}
	for _, item := range v3.Content {
		if item.Label != "newentry" {
			t.Fatalf("expected label 'newentry', got %q", item.Label)
		}
	}
}
