package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func TestSaveIncremental(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	v := New(pub, "test-device")
	v.Add(strings.NewReader("pass1"), "e1")
	v.Add(strings.NewReader("pass2"), "e2")
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	v2.Add(strings.NewReader("pass3"), "e3")
	v2.Add(strings.NewReader("pass4"), "e4")
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	v3, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if len(v3.content) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(v3.content))
	}
}
