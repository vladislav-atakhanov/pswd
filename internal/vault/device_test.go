package vault

import (
	"io"
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
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
	v := New(pub1, "device1")
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
	v2, err := Open(mf, mf.Len(), priv1)
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
	v3, err := Open(mf, mf.Len(), priv2)
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
	v4, err := Open(mf, mf.Len(), priv2)
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
	_, err = Open(mf, mf.Len(), priv1)
	if err == nil {
		t.Fatal("expected error opening with removed device, got nil")
	}
}

func TestAddDeviceBothCanRead(t *testing.T) {
	priv1, pub1, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	priv2, pub2, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}

	// Create vault with device1 and 2 entries
	v := New(pub1, "device1")
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
	v2, err := Open(mf, mf.Len(), priv1)
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

	// Read all entries with priv1
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v3, err := Open(mf, mf.Len(), priv1)
	if err != nil {
		t.Fatalf("open with priv1 after AddDevice: %v", err)
	}
	for id := range v3.Content {
		r, err := v3.Read(mf, id, priv1)
		if err != nil {
			t.Fatalf("read entry with priv1: %v", err)
		}
		_, err = io.ReadAll(r)
		if err != nil {
			t.Fatalf("read all with priv1: %v", err)
		}
	}

	// Read all entries with priv2
	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v4, err := Open(mf, mf.Len(), priv2)
	if err != nil {
		t.Fatalf("open with priv2 after AddDevice: %v", err)
	}
	expectedLabels := map[string]bool{"entry1": true, "entry2": true}
	for id, item := range v4.Content {
		r, err := v4.Read(mf, id, priv2)
		if err != nil {
			t.Fatalf("read entry with priv2: %v", err)
		}
		_, err = io.ReadAll(r)
		if err != nil {
			t.Fatalf("read all with priv2: %v", err)
		}
		delete(expectedLabels, item.Label)
	}
	if len(expectedLabels) > 0 {
		t.Fatalf("missing entries: %v", expectedLabels)
	}
}
