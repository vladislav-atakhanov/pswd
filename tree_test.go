package pswd_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vladislav-atakhanov/pswd"
)

func TestType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pswd-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := tmpDir

	// Create a directory and a password file
	if err := os.Mkdir(filepath.Join(storagePath, "a-dir"), 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	if _, err := os.Create(filepath.Join(storagePath, "a-file.asc")); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	p, err := pswd.NewPswd(storagePath, "")
	if err != nil {
		t.Fatalf("Failed to create pswd: %v", err)
	}

	tests := []struct {
		name         string
		expectedType pswd.Passtype
		expectError  bool
	}{
		{name: "", expectedType: pswd.PassDir, expectError: false},
		{name: "a-dir", expectedType: pswd.PassDir, expectError: false},
		{name: "a-file", expectedType: pswd.PassFile, expectError: false},
		{name: "non-existent", expectedType: pswd.PassUnknown, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := p.Type(tt.name)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if gotType != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, gotType)
			}
		})
	}
}

func TestTree(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pswd-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := tmpDir

	// Create a directory structure
	paths := []string{
		"a/a1.asc",
		"a/a2.asc",
		"b/b1.asc",
		".git/config",
		"c/c1/c11.asc",
	}

	for _, p := range paths {
		path := filepath.Join(storagePath, p)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", p, err)
		}
		if _, err := os.Create(path); err != nil {
			t.Fatalf("Failed to create file %s: %v", p, err)
		}
	}

	p, err := pswd.NewPswd(storagePath, "")
	if err != nil {
		t.Fatalf("Failed to create pswd: %v", err)
	}

	t.Run("tree on root", func(t *testing.T) {
		tree, err := p.Tree("")
		if err != nil {
			t.Fatalf("Tree failed: %v", err)
		}

		expected := &pswd.TreeNode{
			Name:  filepath.Base(storagePath),
			IsDir: true,
			Children: []*pswd.TreeNode{
				{Name: "a", IsDir: true, Children: []*pswd.TreeNode{
					{Name: "a1", IsDir: false},
					{Name: "a2", IsDir: false},
				}},
				{Name: "b", IsDir: true, Children: []*pswd.TreeNode{
					{Name: "b1", IsDir: false},
				}},
				{Name: "c", IsDir: true, Children: []*pswd.TreeNode{
					{Name: "c1", IsDir: true, Children: []*pswd.TreeNode{
						{Name: "c11", IsDir: false},
					}},
				}},
			},
		}

		if !reflect.DeepEqual(tree, expected) {
			t.Errorf("Tree structure mismatch.\nGot:      %#v\nExpected: %#v", tree, expected)
		}
	})

	t.Run("tree on subdirectory", func(t *testing.T) {
		tree, err := p.Tree("a")
		if err != nil {
			t.Fatalf("Tree failed: %v", err)
		}

		expected := &pswd.TreeNode{
			Name:  "a",
			IsDir: true,
			Children: []*pswd.TreeNode{
				{Name: "a1", IsDir: false},
				{Name: "a2", IsDir: false},
			},
		}

		if !reflect.DeepEqual(tree, expected) {
			t.Errorf("Tree structure mismatch.\nGot:      %#v\nExpected: %#v", tree, expected)
		}
	})

	t.Run("tree on file", func(t *testing.T) {
		tree, err := p.Tree("a/a1")
		if err != nil {
			t.Fatalf("Tree failed: %v", err)
		}

		expected := &pswd.TreeNode{
			Name:  "a1",
			IsDir: false,
		}

		if !reflect.DeepEqual(tree, expected) {
			t.Errorf("Tree structure mismatch.\nGot:      %#v\nExpected: %#v", tree, expected)
		}
	})

	t.Run("tree on non-existent path", func(t *testing.T) {
		_, err := p.Tree("non-existent")
		if err == nil {
			t.Errorf("Expected error for non-existent path, got nil")
		}
	})
}
