package vault

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

// Vault Паспорт всего нашего хранилища в RAM
type Vault struct {
	Devices []Device
	Index   []IndexEntry
	Data    []byte

	cursor  int
	dataEnd int
	Full    bool
}

func Open(r io.ReadSeeker, size int) (*Vault, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	version, err := readString(r, 4)
	if err != nil {
		return nil, err
	}
	if version != "VLT1" {
		return nil, fmt.Errorf("Unknown version %s", version)
	}
	devices, err := readDevices(r)
	if err != nil {
		return nil, err
	}
	index, err := readIndexes(r)
	if err != nil {
		return nil, err
	}
	v := new(Vault{Devices: devices, Index: index})
	v.dataEnd = size - v.FooterLength()
	return v, nil
}

const version_length = 4

func (v *Vault) HeaderLength() int {
	res := version_length + 2
	for _, d := range v.Devices {
		res += len(d.Bytes())
	}
	return res
}
func (v *Vault) FooterLength() int {
	res := 4
	for _, d := range v.Index {
		res += len(d.Bytes())
	}
	return res
}

func readString(file io.ReadSeeker, length int) (string, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(file, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (v *Vault) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "Devices (%d):\n", len(v.Devices))
	for _, d := range v.Devices {
		key := d.PublicKey()
		fmt.Fprintf(b, "\t%s %s\n", d.Name(), base64.URLEncoding.EncodeToString(key[:]))
	}
	fmt.Fprintf(b, "Index (%d):\n", len(v.Index))
	for _, i := range v.Index {
		fmt.Fprintf(b, "\t%s (%d:%d)\n", i.Label(), v.HeaderLength()+i.Start(), i.Length())
	}
	return b.String()
}

func (v *Vault) ReadRange(r io.ReaderAt, out io.Writer, privateKey [32]byte, from, lenght int) error {
	publicKey, err := crypto.PublicKeyFromPrivate(privateKey)
	if err != nil {
		return err
	}

	index := -1
	for i, d := range v.Devices {
		p := d.PublicKey()
		if bytes.Equal(p[:], publicKey[:]) {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("Access denied")
	}
	if err := crypto.DecryptStream(out, io.NewSectionReader(r, int64(from), int64(lenght)), privateKey, len(v.Devices), index); err != nil {
		return err
	}
	return nil
}
