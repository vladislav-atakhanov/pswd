package vault

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

var (
	ErrAccessDenied   = errors.New("access denied")
	ErrNotFound       = errors.New("password not found")
	ErrDeviceExists   = errors.New("device already in vault")
	ErrDeviceNotFound = errors.New("device not found")
	ErrInvalidLabel   = errors.New("invalid label")
)

type Item struct {
	Label      string
	content    io.Reader
	start      int
	length     int
	LastUpdate uint64
}
type contentKey = uuid.V4

// Vault holds the entire vault state in memory
type Vault struct {
	Devices []Device
	content map[contentKey]Item

	dataEnd int
	Full    bool

	orphanedSpans []Span
}

func New(publicKey [32]byte, name string) *Vault {
	return &Vault{
		Devices: []Device{newDevice(publicKey, name)},
		content: make(map[contentKey]Item),
		Full:    true,
	}
}
func Open(r io.ReadSeeker, size int, privateKey [32]byte) (*Vault, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	version, err := readString(r, 4)
	if err != nil {
		return nil, err
	}
	if version != "VLT1" {
		return nil, fmt.Errorf("unknown version: %s", version)
	}
	devices, err := readDevices(r)
	if err != nil {
		return nil, err
	}
	v := new(Vault{Devices: devices,
		content: make(map[contentKey]Item),
	})
	if err := v.read(r, size, privateKey); err != nil {
		return nil, err
	}
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

func readString(file io.ReadSeeker, length int) (string, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(file, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (v *Vault) Print(b io.Writer) error {
	if _, err := fmt.Fprintf(b, "Devices (%d):\n", len(v.Devices)); err != nil {
		return err
	}
	for _, d := range v.Devices {
		key := d.PublicKey()
		if _, err := fmt.Fprintf(b, "\t%s %s\n", d.Name(), base64.URLEncoding.EncodeToString(key[:])); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(b, "Passwords (%d):\n", len(v.content)); err != nil {
		return err
	}
	for id, i := range v.content {
		if _, err := fmt.Fprintf(b, "\t%s | %s (%d:%d)\n", id.String(), i.Label, i.start, i.length); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vault) decrypt(r io.Reader, out io.Writer, privateKey [32]byte) error {
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
		return ErrAccessDenied
	}
	if err := crypto.DecryptStream(out, r, privateKey, len(v.Devices), index); err != nil {
		return err
	}
	return nil
}

func (v *Vault) Read(r io.ReaderAt, id contentKey, privateKey [32]byte) (io.Reader, error) {
	item, ok := v.content[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id.String())
	}
	if item.content != nil {
		return item.content, nil
	}
	var buffer bytes.Buffer
	if err := v.decrypt(io.NewSectionReader(r, int64(v.HeaderLength())+int64(item.start), int64(item.length)), &buffer, privateKey); err != nil {
		return nil, err
	}
	return &buffer, nil
}
