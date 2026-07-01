package vault

import "io"

func (v *Vault) Compact(r io.ReaderAt, privateKey [32]byte) error {
	if err := v.loadContent(r, privateKey); err != nil {
		return err
	}
	v.Full = true
	v.orphanedSpans = nil
	return nil
}
