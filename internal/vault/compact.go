package vault

import (
	"bytes"
	"io"
)

func (v *Vault) Compact(r io.ReaderAt, privateKey [32]byte) error {
	for id, item := range v.Content {
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
		v.Content[id] = item
	}

	v.Full = true
	v.orphanedSpans = nil
	return nil
}
