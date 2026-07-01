package main

import (
	"os"
	"testing"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"y\n", true},
		{"yes\n", true},
		{"Y\n", true},
		{"YES\n", true},
		{"n\n", false},
		{"no\n", false},
		{"N\n", false},
		{"\n", false},
		{"foo\n", false},
	}
	for _, tt := range tests {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(tt.input))
		w.Close()

		oldStdin := os.Stdin
		os.Stdin = r

		got := confirm("prompt: ")

		os.Stdin = oldStdin
		r.Close()

		if got != tt.want {
			t.Errorf("confirm(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
