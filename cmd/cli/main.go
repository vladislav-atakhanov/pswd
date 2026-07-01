package main

import (
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
	v := vault.New()
	priv, err := readKey(private)
	v.AddDevice(pub, "pc", nil, priv)
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
	defer file.Close()
	priv, err := readKey(private)
	if err != nil {
		panic(err)
	}
	v, err := vault.Open(file, size, priv)
	if err != nil {
		if v, err = initStorage(); err != nil {
			panic(err)
		}
	}
	defer func() {
		v.Print(os.Stdout)
		if err := v.Save(file); err != nil {
			panic(err)
		}
	}()
	// content := must(v.Read(file, must(uuid.UUIDv4FromString("9fcbb1c4-3dbe-49e3-bd6d-945a359ea6a8")), priv))
	// fmt.Println(content)
	// v.Add([]byte("yout"), "youtube")
	// v.Add([]byte("pass"), "github")
	for id, i := range v.Content {
		content := must(v.Read(file, id, priv))
		fmt.Println(i.Label, "-", content)
	}

	// priv, err := readKey(private)
	// if err != nil {
	// 	panic(err)
	// }
	// var buf bytes.Buffer
	// if err := v.ReadRange(file, &buf, priv, 158, 115); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(buf.String())
}
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
