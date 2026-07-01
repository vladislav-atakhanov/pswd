package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

func TestRename(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

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

	// Open, rename
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	var id uuid.V4
	for id = range v2.content {
		break
	}
	if err := v2.Rename(id, "renamed-entry"); err != nil {
		t.Fatal(err)
	}

	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Reopen — verify new label
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after rename: %v", err)
	}
	if len(v3.content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.content))
	}
	for _, item := range v3.content {
		if item.Label != "renamed-entry" {
			t.Fatalf("expected label 'renamed-entry', got %q", item.Label)
		}
	}
}
