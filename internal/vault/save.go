package vault

import (
	"encoding/binary"
	"io"
)

func (v *Vault) SaveFull(w io.WriteSeeker) error {
	if _, err := w.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := v.saveHeader(w); err != nil {
		return nil
	}
	if _, err := w.Write(v.Data); err != nil {
		return err
	}
	if err := v.saveIndex(w); err != nil {
		return err
	}
	return nil
}
func (v *Vault) Save(w io.WriteSeeker) error {
	if v.Full {
		return v.SaveFull(w)
	}
	if _, err := w.Seek(int64(v.dataEnd), io.SeekStart); err != nil {
		return err
	}
	if _, err := w.Write(v.Data); err != nil {
		return err
	}
	if err := v.saveIndex(w); err != nil {
		return err
	}
	return nil
}

func (v *Vault) saveDevices(file io.Writer) error {
	if err := binary.Write(file, binary.BigEndian, uint16(len(v.Devices))); err != nil {
		return err
	}
	for _, dev := range v.Devices {
		if _, err := file.Write(dev.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vault) saveIndex(file io.Writer) error {
	for _, entry := range v.Index {
		if _, err := file.Write(entry.Bytes()); err != nil {
			return err
		}
	}
	if err := binary.Write(file, binary.BigEndian, uint32(len(v.Index))); err != nil {
		return err
	}
	return nil
}

func (v *Vault) saveHeader(file io.Writer) error {
	if _, err := file.Write([]byte("VLT1")); err != nil {
		return err
	}
	if err := v.saveDevices(file); err != nil {
		return err
	}
	return nil
}
