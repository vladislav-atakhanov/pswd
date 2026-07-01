package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

var exportCmd = &cobra.Command{
	Use:   "export <keyfile> <name>",
	Short: "Export the public key with a device name",
	Long:  "Decrypt the private key, derive the public key, and output a base64-encoded blob of the public key, name length, and name.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyfile := args[0]
		name := args[1]

		data, err := os.ReadFile(keyfile)
		if err != nil {
			return fmt.Errorf("read key file: %w", err)
		}

		password, err := readPassword("Enter master password: ")
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}
		defer mem.ZeroBytes(password)

		priv, err := crypto.DecryptPrivateKey(data, password)
		if err != nil {
			return fmt.Errorf("decrypt private key: %w", err)
		}
		defer mem.ZeroArray32(&priv)

		pub, err := crypto.PublicKeyFromPrivate(priv)
		if err != nil {
			return fmt.Errorf("derive public key: %w", err)
		}

		nameBytes := []byte(name)
		var length [2]byte
		binary.BigEndian.PutUint16(length[:], uint16(len(nameBytes)))

		buf := make([]byte, 0, 32+2+len(nameBytes))
		buf = append(buf, pub[:]...)
		buf = append(buf, length[:]...)
		buf = append(buf, nameBytes...)

		fmt.Println(base64.URLEncoding.EncodeToString(buf))

		return nil
	},
}
