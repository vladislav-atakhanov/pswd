package vault

import (
	"strings"
	"testing"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func FuzzOpen(f *testing.F) {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		f.Fatal(err)
	}

	mf := &mem.MemoryFile{}
	v := New(pub, "test-device")
	if err := v.Add(strings.NewReader("password"), "entry1"); err != nil {
		f.Fatal(err)
	}
	if err := v.Save(mf); err != nil {
		f.Fatal(err)
	}
	f.Add(mf.Bytes())

	mf2 := &mem.MemoryFile{}
	v2 := New(pub, "test-device")
	if err := v2.Save(mf2); err != nil {
		f.Fatal(err)
	}
	f.Add(mf2.Bytes())

	f.Add([]byte{})
	f.Add([]byte("Hello, World!"))

	f.Fuzz(func(t *testing.T, data []byte) {
		mf := &mem.MemoryFile{}
		mf.Write(data)
		Open(mf, mf.Len(), priv)
	})
}
