package vault

import (
	"errors"
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

func TestValidateLabel(t *testing.T) {
	tests := []struct {
		label string
		err   bool
	}{
		{"valid-label", false},
		{"a", false},
		{"", true},
	}
	for _, tc := range tests {
		err := validateLabel(tc.label)
		if tc.err && err == nil {
			t.Errorf("expected error for label %q, got nil", tc.label)
		}
		if !tc.err && err != nil {
			t.Errorf("unexpected error for label %q: %v", tc.label, err)
		}
	}
}

func TestValidateLabelTooLong(t *testing.T) {
	label := string(make([]byte, MaxLabelLength+1))
	err := validateLabel(label)
	if err == nil {
		t.Fatal("expected error for label exceeding MaxLabelLength, got nil")
	}
	if !errors.Is(err, ErrInvalidLabel) {
		t.Fatalf("expected ErrInvalidLabel, got %v", err)
	}
}

func TestAddEmptyLabel(t *testing.T) {
	_, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	v := New(pub, "test-device")
	err = v.Add(strings.NewReader("password"), "")
	if err == nil {
		t.Fatal("expected error for empty label, got nil")
	}
	if !errors.Is(err, ErrInvalidLabel) {
		t.Fatalf("expected ErrInvalidLabel, got %v", err)
	}
}

func TestRenameEmptyLabel(t *testing.T) {
	_, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	v := New(pub, "test-device")
	if err := v.Add(strings.NewReader("password"), "entry1"); err != nil {
		t.Fatal(err)
	}
	var id contentKey
	for id = range v.content {
		break
	}
	err = v.Rename(id, "")
	if err == nil {
		t.Fatal("expected error for empty label, got nil")
	}
	if !errors.Is(err, ErrInvalidLabel) {
		t.Fatalf("expected ErrInvalidLabel, got %v", err)
	}
}

func TestRemoveNotFound(t *testing.T) {
	_, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	v := New(pub, "test-device")
	var id uuid.V4
	err = v.Remove(id)
	if err == nil {
		t.Fatal("expected error for removing non-existent ID, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRemoveNotFoundAfterSave(t *testing.T) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	v := New(pub, "test-device")
	if err := v.Add(strings.NewReader("password1"), "entry1"); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		t.Fatal(err)
	}

	if _, err := mf.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	v2, err := Open(mf, mf.Len(), priv)
	if err != nil {
		t.Fatal(err)
	}

	var fakeID uuid.V4
	err = v2.Remove(fakeID)
	if err == nil {
		t.Fatal("expected error for removing non-existent ID, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
