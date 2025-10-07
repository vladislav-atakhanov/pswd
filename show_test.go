package pswd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vladislav-atakhanov/pswd"
	"github.com/vladislav-atakhanov/pswd/pkg/keys"
)

func TestShow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pswd-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "storage")
	keysPath := filepath.Join(tmpDir, "keys")

	if err := os.MkdirAll(storagePath, 0700); err != nil {
		t.Fatalf("Failed to create storage path: %v", err)
	}
	if err := os.MkdirAll(keysPath, 0700); err != nil {
		t.Fatalf("Failed to create keys path: %v", err)
	}

	keyID := "testuser@example.com"
	masterPassword := "masterpassword"
	priv, pub, err := keys.Generate(masterPassword)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	ks, err := keys.NewKeyStore(keysPath)
	if err != nil {
		t.Fatalf("Failed to create keystore: %v", err)
	}

	if _, err := ks.Save(keyID, priv, pub); err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	if err := os.WriteFile(filepath.Join(storagePath, ".key-id"), []byte(keyID), 0644); err != nil {
		t.Fatalf("Failed to write key id file: %v", err)
	}

	p, err := pswd.NewPswd(storagePath, keysPath)
	if err != nil {
		t.Fatalf("Failed to create pswd: %v", err)
	}

	passwordName := "mypass"
	passwordValue := "mypassword"

	if _, err := p.Insert(passwordName, passwordValue); err != nil {
		t.Fatalf("Failed to insert password: %v", err)
	}

	t.Run("show password with correct master password", func(t *testing.T) {
		decrypted, err := p.Show(passwordName, masterPassword)
		if err != nil {
			t.Fatalf("Show failed: %v", err)
		}
		if decrypted != passwordValue {
			t.Errorf("Expected password %q, got %q", passwordValue, decrypted)
		}
	})

	t.Run("show password with incorrect master password", func(t *testing.T) {
		_, err := p.Show(passwordName, "wrongpassword")
		if err == nil {
			t.Errorf("Expected error for incorrect master password, got nil")
		}
	})

	t.Run("show non-existent password", func(t *testing.T) {
		_, err := p.Show("nonexistent", masterPassword)
		if err == nil {
			t.Errorf("Expected error for non-existent password, got nil")
		}
	})
}

func TestShowLazy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pswd-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "storage")
	keysPath := filepath.Join(tmpDir, "keys")

	if err := os.MkdirAll(storagePath, 0700); err != nil {
		t.Fatalf("Failed to create storage path: %v", err)
	}
	if err := os.MkdirAll(keysPath, 0700); err != nil {
		t.Fatalf("Failed to create keys path: %v", err)
	}

	keyID := "testuser@example.com"
	masterPassword := "masterpassword"
	priv, pub, err := keys.Generate(masterPassword)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	ks, err := keys.NewKeyStore(keysPath)
	if err != nil {
		t.Fatalf("Failed to create keystore: %v", err)
	}

	if _, err := ks.Save(keyID, priv, pub); err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	if err := os.WriteFile(filepath.Join(storagePath, ".key-id"), []byte(keyID), 0644); err != nil {
		t.Fatalf("Failed to write key id file: %v", err)
	}

	p, err := pswd.NewPswd(storagePath, keysPath)
	if err != nil {
		t.Fatalf("Failed to create pswd: %v", err)
	}

	passwordName := "mypass"
	passwordValue := "mypassword"

	if _, err := p.Insert(passwordName, passwordValue); err != nil {
		t.Fatalf("Failed to insert password: %v", err)
	}

	t.Run("show password with correct master password", func(t *testing.T) {
		decrypted, err := p.ShowLazy(passwordName, func(key string) (string, error) {
			return masterPassword, nil
		})
		if err != nil {
			t.Fatalf("ShowLazy failed: %v", err)
		}
		if decrypted != passwordValue {
			t.Errorf("Expected password %q, got %q", passwordValue, decrypted)
		}
	})

	t.Run("show password with incorrect master password", func(t *testing.T) {
		_, err := p.ShowLazy(passwordName, func(key string) (string, error) {
			return "wrongpassword", nil
		})
		if err == nil {
			t.Errorf("Expected error for incorrect master password, got nil")
		}
	})

	t.Run("show non-existent password", func(t *testing.T) {
		_, err := p.ShowLazy("nonexistent", func(key string) (string, error) {
			return masterPassword, nil
		})
		if err == nil {
			t.Errorf("Expected error for non-existent password, got nil")
		}
	})
}
