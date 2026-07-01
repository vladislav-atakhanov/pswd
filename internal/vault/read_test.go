package vault

import (
	"errors"
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func TestReadWrongKey(t *testing.T) {
	priv1, pub1, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	priv2, _, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	v := New(pub1, "device1")
	if err := v.Add(strings.NewReader("password"), "entry"); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	v2, err := Open(mf, mf.Len(), priv1)
	if err != nil {
		t.Fatalf("open with priv1: %v", err)
	}

	var id contentKey
	for id = range v2.content {
		break
	}

	_, err = v2.Read(mf, id, priv2)
	if err == nil {
		t.Fatal("expected error reading with wrong private key, got nil")
	}
	if !errors.Is(err, ErrAccessDenied) {
		t.Fatalf("expected ErrAccessDenied, got %v", err)
	}
}
