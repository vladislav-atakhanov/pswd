package mem

import (
	"io"
	"testing"
)

func TestMemoryFileWriteRead(t *testing.T) {
	var mf MemoryFile

	n, err := mf.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if n != 11 {
		t.Fatalf("expected 11 bytes written, got %d", n)
	}

	if _, err := mf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 11)
	n, err = mf.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if n != 11 {
		t.Fatalf("expected 11 bytes read, got %d", n)
	}
	if string(buf) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(buf))
	}
}

func TestMemoryFileSeekEnd(t *testing.T) {
	var mf MemoryFile
	mf.Write([]byte("hello"))

	pos, err := mf.Seek(0, io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 5 {
		t.Fatalf("expected position 5, got %d", pos)
	}

	mf.Write([]byte(" world"))
	if _, err := mf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	buf, err := io.ReadAll(&mf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(buf))
	}
}

func TestMemoryFileTruncate(t *testing.T) {
	var mf MemoryFile
	mf.Write([]byte("hello world"))

	if err := mf.Truncate(5); err != nil {
		t.Fatal(err)
	}
	if mf.Len() != 5 {
		t.Fatalf("expected length 5, got %d", mf.Len())
	}

	if _, err := mf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 5)
	mf.Read(buf)
	if string(buf) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(buf))
	}

	mf.Write([]byte(" there"))
	if _, err := mf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	buf, err := io.ReadAll(&mf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "hello there" {
		t.Fatalf("expected 'hello there', got %q", string(buf))
	}
}

func TestMemoryFileReadAt(t *testing.T) {
	var mf MemoryFile
	mf.Write([]byte("hello world"))

	buf := make([]byte, 5)
	n, err := mf.ReadAt(buf, 6)
	if err != nil {
		t.Fatalf("ReadAt: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5 bytes, got %d", n)
	}
	if string(buf) != "world" {
		t.Fatalf("expected 'world', got %q", string(buf))
	}
}

func TestMemoryFileGrowth(t *testing.T) {
	var mf MemoryFile

	chunk := make([]byte, 100)
	for i := 0; i < 10; i++ {
		mf.Write(chunk)
	}

	if mf.Len() < 1000 {
		t.Fatalf("expected at least 1000 bytes after 10 writes, got %d", mf.Len())
	}
}
