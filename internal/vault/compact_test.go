package vault

import (
	"io"
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func TestCompact(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with 3 entries and save
	v := New(pub, "test")
	if err := v.Add(strings.NewReader("password1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("password2"), "entry2"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("password3"), "entry3"); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	// Open, remove one entry, save
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatal(err)
	}
	var idToRemove contentKey
	for id := range v2.content {
		idToRemove = id
		break
	}
	if err := v2.Remove(idToRemove); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}
	sizeAfterRemove := mf.Len()

	// Compact and save
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatal(err)
	}
	if err := v3.Compact(mf, priv); err != nil {
		t.Fatal(err)
	}
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := v3.Save(mf); err != nil {
		t.Fatal(err)
	}
	sizeAfterCompact := mf.Len()

	// Compact must shrink the file (no zeroed orphaned blobs)
	if sizeAfterCompact >= sizeAfterRemove {
		t.Fatalf("expected compact to shrink file: %d >= %d", sizeAfterCompact, sizeAfterRemove)
	}

	// Reopen and verify 2 entries remain, all readable
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v4, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen after compact: %v", err)
	}
	if len(v4.content) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(v4.content))
	}
	for id := range v4.content {
		r, err := v4.Read(mf, id, priv)
		if err != nil {
			t.Fatalf("read after compact: %v", err)
		}
		if _, err := io.ReadAll(r); err != nil {
			t.Fatalf("read content after compact: %v", err)
		}
	}
}
