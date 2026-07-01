package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

func TestList(t *testing.T) {
	_, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	v := New()
	if err := v.InitDevice(pub, "test"); err != nil {
		t.Fatal(err)
	}

	if err := v.Add(strings.NewReader("pass1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Add(strings.NewReader("pass2"), "entry2"); err != nil {
		t.Fatal(err)
	}

	entries := v.List()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	labels := make(map[string]bool)
	for _, e := range entries {
		labels[e.Label] = true
		if e.ID.String() == "00000000-0000-0000-0000-000000000000" {
			t.Fatal("expected non-zero ID")
		}
		if e.LastUpdate == 0 {
			t.Fatal("expected non-zero LastUpdate")
		}
	}
	if !labels["entry1"] {
		t.Fatal("missing entry1")
	}
	if !labels["entry2"] {
		t.Fatal("missing entry2")
	}
}
