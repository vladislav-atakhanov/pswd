package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

func TestRemoveOneOfTwo(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// First session: create vault with 2 entries
	v := New(pub, "test-device")
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
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	if len(v2.content) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(v2.content))
	}

	var id uuid.V4
	for id = range v2.content {
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
	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("second open (after remove+save): %v", err)
	}
	if len(v3.content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.content))
	}
}

func TestRemoveAll(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with 2 entries
	v := New(pub, "test-device")
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
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	for id := range v2.content {
		if err := v2.Remove(id); err != nil {
			t.Fatal(err)
		}
	}
	if len(v2.content) != 0 {
		t.Fatalf("expected 0 entries after remove all, got %d", len(v2.content))
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
	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after remove all: %v", err)
	}
	if len(v3.content) != 0 {
		t.Fatalf("expected empty vault, got %d", len(v3.content))
	}
}

func TestRemoveThenAdd(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with 1 entry
	v := New(pub, "test-device")
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
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	for id := range v2.content {
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
	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after remove+add: %v", err)
	}
	if len(v3.content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.content))
	}
	for _, item := range v3.content {
		if item.Label != "newentry" {
			t.Fatalf("expected label 'newentry', got %q", item.Label)
		}
	}
}
