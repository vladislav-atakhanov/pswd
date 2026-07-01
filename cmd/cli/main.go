package main

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

const private = "private.age"
const public = "public.age"

func readPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return password, err
}

func generateKeys() {
	priv, pub, err := crypto.GenerateKeys()
	if err != nil {
		panic(err)
	}
	password, err := readPassword("Enter master password: ")
	if err != nil {
		panic(err)
	}
	defer mem.ZeroBytes(password)
	encrypted, err := crypto.EncryptPrivateKey(priv, password)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(private, encrypted, 0600); err != nil {
		panic(err)
	}
	if err := os.WriteFile(public, pub[:], 0644); err != nil {
		panic(err)
	}
}

func readPublicKey(filename string) ([32]byte, error) {
	pub, err := os.ReadFile(filename)
	if err != nil {
		return [32]byte{}, err
	}
	var p [32]byte
	copy(p[:], pub)
	return p, nil
}

func readPrivateKey(filename string) ([32]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return [32]byte{}, err
	}
	password, err := readPassword("Enter master password: ")
	if err != nil {
		return [32]byte{}, err
	}
	defer mem.ZeroBytes(password)
	return crypto.DecryptPrivateKey(data, password)
}

func initStorage() (*vault.Vault, error) {
	pub, err := readPublicKey(public)
	if err != nil {
		return nil, err
	}
	v := vault.New()
	if err := v.InitDevice(pub, "pc"); err != nil {
		return nil, err
	}
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
	priv, err := readPrivateKey(private)
	if err != nil {
		panic(err)
	}
	defer mem.ZeroArray32(&priv)
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
	// v.Add(strings.NewReader("mail"), "mail.ru")
	// v.Add(strings.NewReader("pass"), "github")
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
