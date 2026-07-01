package vault

import (
	"encoding/binary"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

type Span struct {
	Start  int
	Length int
}

func (v *Vault) saveBody(w io.Writer) error {
	keys := make([][32]byte, len(v.Devices))
	for i, d := range v.Devices {
		keys[i] = d.PublicKey()
	}

	cursor := v.dataEnd - v.HeaderLength()
	spans := make(map[contentKey]Span)
	for id, i := range v.Content {
		if i.content == nil {
			spans[id] = Span{Length: i.length, Start: i.start}
			continue
		}

		n, err := crypto.EncryptStream(w, i.content, keys)
		if err != nil {
			return err
		}
		spans[id] = Span{Length: n, Start: cursor}
		cursor += n
	}
	read, write := io.Pipe()
	defer read.Close()
	go func(w io.WriteCloser) error {
		defer w.Close()
		if err := binary.Write(w, binary.BigEndian, uint32(len(v.Content))); err != nil {
			return err
		}
		for id, i := range v.Content {
			if _, err := w.Write(id[:]); err != nil {
				return err
			}
			if err := binary.Write(w, binary.BigEndian, i.LastUpdate); err != nil {
				return err
			}
			span := spans[id]
			if err := binary.Write(w, binary.BigEndian, uint32(span.Start)); err != nil {
				return err
			}
			if err := binary.Write(w, binary.BigEndian, uint32(span.Length)); err != nil {
				return err
			}
			labelBytes := []byte(i.Label)
			if err := binary.Write(w, binary.BigEndian, uint16(len(labelBytes))); err != nil {
				return err
			}
			if _, err := w.Write(labelBytes); err != nil {
				return err
			}
		}
		return nil
	}(write)

	n, err := crypto.EncryptStream(w, read, keys)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint32(n)); err != nil {
		return err
	}
	return nil
}

type Writer interface {
	io.Seeker
	io.Writer
}

func (v *Vault) Save(w Writer) error {
	if v.Full {
		if _, err := w.Seek(0, io.SeekStart); err != nil {
			return err
		}
		if err := v.saveHeader(w); err != nil {
			return err
		}
		v.dataEnd = v.HeaderLength()
	} else {
		if _, err := w.Seek(int64(v.dataEnd), io.SeekStart); err != nil {
			return err
		}
	}
	if err := v.saveBody(w); err != nil {
		return err
	}

	return nil
}

func (v *Vault) saveHeader(w io.Writer) error {
	if _, err := w.Write([]byte("VLT1")); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(len(v.Devices))); err != nil {
		return err
	}
	for _, dev := range v.Devices {
		if _, err := w.Write(dev.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
