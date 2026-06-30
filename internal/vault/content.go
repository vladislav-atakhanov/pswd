package vault

import (
	"bytes"
	"crypto/rand"
	"time"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

func NewUUID() ([16]byte, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return uuid, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid, nil
}

func (v *Vault) Add(plain []byte, label string) error {
	id, err := NewUUID()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	keys := make([][32]byte, len(v.Devices))
	for i, d := range v.Devices {
		keys[i] = d.PublicKey()
	}
	crypto.EncryptStream(&buffer, bytes.NewReader(plain), keys)

	content := buffer.Bytes()
	v.Data = append(v.Data, content...)
	v.Index = append(v.Index, newIndex(
		label,
		len(content),
		v.cursor,
		uint64(time.Now().Unix()),
		id,
	))
	v.cursor += len(content)
	return nil
}
