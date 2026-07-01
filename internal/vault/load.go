package vault

import (
	"bytes"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func (v *Vault) loadContent(r io.ReaderAt, privateKey [32]byte) error {
	defer mem.ZeroArray32(&privateKey)

	for id, item := range v.content {
		if item.content != nil {
			continue
		}
		b, err := v.Read(r, id, privateKey)
		if err != nil {
			return err
		}
		pass, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		item.content = bytes.NewReader(pass)
		v.content[id] = item
	}
	return nil
}
