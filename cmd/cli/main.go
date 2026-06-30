package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

const private = "private.age"
const public = "public.age"

func generateKeys() {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		panic(err)
	}
	os.WriteFile(public, pub[:], 0667)
	os.WriteFile(private, priv[:], 0667)
}

func readKey(filename string) ([32]byte, error) {
	pub, err := os.ReadFile(filename)
	if err != nil {
		return [32]byte{}, err
	}
	var p [32]byte
	copy(p[:], pub)
	return p, nil
}

func initStorage() (*vault.Vault, error) {
	pub, err := readKey(public)
	if err != nil {
		return nil, err
	}
	v := &vault.Vault{Full: true}
	v.AddDevice(pub, "pc")
	return v, nil
}
func open(name string) (*os.File, int, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0667)
	if err == nil {
		stat, err := file.Stat()
		if err != nil {
			return nil, 0, err
		}

		return file, int(stat.Size()), nil
	}
	file, err = os.Create(name)
	if err != nil {
		return nil, 0, err
	}
	return file, 0, nil
}

func main() {

	file, size, err := open("test.pswd")
	if err != nil {
		panic(err)
	}
	v, err := vault.Open(file, size)
	if err != nil {
		if v, err = initStorage(); err != nil {
			panic(err)
		}
	}
	defer v.Save(file)
	// v.Add([]byte("password"), "github")
	// v.Add([]byte("napojlb"), "youtube")
	fmt.Println(v.String())

	priv, err := readKey(private)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := v.ReadRange(file, &buf, priv, 158, 115); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
