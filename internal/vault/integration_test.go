package vault

import (
	"io"
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func TestCrossDeviceRead(t *testing.T) {
	priv1, pub1, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	priv2, pub2, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// dev1 creates vault with one entry
	v := New(pub1, "dev1")
	if err := v.Add(strings.NewReader("pass1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	// dev1 opens, adds dev2, saves
	v2, err := Open(mf, mf.Len(), priv1)
	if err != nil {
		t.Fatalf("open with priv1: %v", err)
	}
	if err := v2.AddDevice(pub2, "dev2", mf, priv1); err != nil {
		t.Fatal(err)
	}
	if err := v2.Save(mf); err != nil {
		t.Fatal(err)
	}

	// dev2 opens and reads entry1 created by dev1
	v3, err := Open(mf, mf.Len(), priv2)
	if err != nil {
		t.Fatalf("open with priv2: %v", err)
	}
	if len(v3.content) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(v3.content))
	}
	for id, item := range v3.content {
		if item.Label != "entry1" {
			t.Fatalf("expected label 'entry1', got %q", item.Label)
		}
		r, err := v3.Read(mf, id, priv2)
		if err != nil {
			t.Fatalf("read with priv2: %v", err)
		}
		b, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("read all: %v", err)
		}
		if string(b) != "pass1" {
			t.Fatalf("expected 'pass1', got %q", string(b))
		}
	}

	// dev2 adds entry2, saves
	v3.Add(strings.NewReader("pass2"), "entry2")
	if err := v3.Save(mf); err != nil {
		t.Fatal(err)
	}

	// dev1 opens and reads both entries
	v4, err := Open(mf, mf.Len(), priv1)
	if err != nil {
		t.Fatalf("reopen with priv1: %v", err)
	}
	if len(v4.content) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(v4.content))
	}
	got := map[string]string{}
	for id, item := range v4.content {
		r, err := v4.Read(mf, id, priv1)
		if err != nil {
			t.Fatalf("read with priv1: %v", err)
		}
		b, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("read all: %v", err)
		}
		got[item.Label] = string(b)
	}
	if got["entry1"] != "pass1" {
		t.Fatalf("expected entry1='pass1', got %q", got["entry1"])
	}
	if got["entry2"] != "pass2" {
		t.Fatalf("expected entry2='pass2', got %q", got["entry2"])
	}
}
