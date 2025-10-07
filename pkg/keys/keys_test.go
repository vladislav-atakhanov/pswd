package keys_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vladislav-atakhanov/pswd/pkg/keys"
)

func TestKeyStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "keys-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ks, err := keys.NewKeyStore(tmpDir)
	if err != nil {
		t.Fatalf("NewKeyStore failed: %v", err)
	}

	keyID := "test@example.com"
	password := "testpassword"

	t.Run("Generate and Save", func(t *testing.T) {
		if ks.Exists(keyID) {
			t.Fatal("Key should not exist before saving")
		}

		priv, pub, err := keys.Generate(password)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if _, err := ks.Save(keyID, priv, pub); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		if !ks.Exists(keyID) {
			t.Fatal("Key should exist after saving")
		}

		// Check if files were created
		if _, err := os.Stat(filepath.Join(tmpDir, keyID, "private.asc")); err != nil {
			t.Errorf("private.asc not found: %v", err)
		}
		if _, err := os.Stat(filepath.Join(tmpDir, keyID, "public.asc")); err != nil {
			t.Errorf("public.asc not found: %v", err)
		}
	})

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		plaintext := []byte("my secret message")

		encrypted, err := ks.Encrypt(keyID, plaintext)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		t.Run("with correct password", func(t *testing.T) {
			decrypted, err := ks.Decrypt(keyID, password, encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if string(decrypted) != string(plaintext) {
				t.Errorf("Decrypted text does not match original. Got %q, want %q", decrypted, plaintext)
			}
		})

		t.Run("with incorrect password", func(t *testing.T) {
			_, err := ks.Decrypt(keyID, "wrongpassword", encrypted)
			if err == nil {
				t.Fatal("Decrypt should fail with incorrect password")
			}
		})
	})

	t.Run("Encrypt and DecryptLazy", func(t *testing.T) {
		plaintext := []byte("my secret message lazy")

		encrypted, err := ks.Encrypt(keyID, plaintext)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		t.Run("with correct password", func(t *testing.T) {
			decrypted, err := ks.DecryptLazy(keyID, func() (string, error) { return password, nil }, encrypted)
			if err != nil {
				t.Fatalf("DecryptLazy failed: %v", err)
			}

			if string(decrypted) != string(plaintext) {
				t.Errorf("Decrypted text does not match original. Got %q, want %q", decrypted, plaintext)
			}
		})

		t.Run("with incorrect password", func(t *testing.T) {
			_, err := ks.DecryptLazy(keyID, func() (string, error) { return "wrongpassword", nil }, encrypted)
			if err == nil {
				t.Fatal("DecryptLazy should fail with incorrect password")
			}
		})
	})
}

func TestNewKeyStoreDefaultPath(t *testing.T) {
	// This test is a bit tricky as it depends on the user's home directory.
	// We are mostly checking that it doesn't return an error.
	ks, err := keys.NewKeyStore("")
	if err != nil {
		t.Fatalf("NewKeyStore with empty path failed: %v", err)
	}
	if ks == nil {
		t.Fatal("NewKeyStore with empty path returned nil")
	}
}
