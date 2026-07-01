package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func TestNew(t *testing.T) {
	_, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	v := New(pub, "test-device")
	if len(v.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(v.Devices))
	}
	if v.Devices[0].Name() != "test-device" {
		t.Fatalf("expected device name 'test-device', got %q", v.Devices[0].Name())
	}
	if !v.Full {
		t.Fatal("expected new vault to be marked Full")
	}
}

func TestOpenInvalidVersion(t *testing.T) {
	priv, _, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	mf.Write([]byte("BAD!"))
	mf.Write(make([]byte, 100))

	_, err = Open(mf, mf.Len(), priv)
	if err == nil {
		t.Fatal("expected error for invalid version, got nil")
	}
}

func TestOpenEmptyFile(t *testing.T) {
	priv, _, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	_, err = Open(mf, 0, priv)
	if err == nil {
		t.Fatal("expected error for empty file, got nil")
	}
}

func TestOpenTruncatedIndex(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	v := New(pub, "test-device")
	if err := v.Add(strings.NewReader("password"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	data := mf.Bytes()
	truncated := data[:len(data)-4]

	mf2 := &mem.MemoryFile{}
	mf2.Write(truncated)

	_, err = Open(mf2, mf2.Len(), priv)
	if err == nil {
		t.Fatal("expected error for truncated index, got nil")
	}
}
