package vault

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
)

type Span struct {
	Start  int
	Length int
}

func (v *Vault) saveBody(w io.Writer) error {
	keys := make([][32]byte, len(v.devices))
	for i, d := range v.devices {
		keys[i] = d.PublicKey()
	}

	cursor := v.dataEnd - v.HeaderLength()
	spans := make(map[contentKey]Span)
	for id, i := range v.content {
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
	var indexBuf bytes.Buffer
	if err := binary.Write(&indexBuf, binary.BigEndian, uint32(len(v.content))); err != nil {
		return err
	}
	for id, i := range v.content {
		if _, err := indexBuf.Write(id[:]); err != nil {
			return err
		}
		if err := binary.Write(&indexBuf, binary.BigEndian, i.LastUpdate); err != nil {
			return err
		}
		span := spans[id]
		if err := binary.Write(&indexBuf, binary.BigEndian, uint32(span.Start)); err != nil {
			return err
		}
		if err := binary.Write(&indexBuf, binary.BigEndian, uint32(span.Length)); err != nil {
			return err
		}
		labelBytes := []byte(i.Label)
		if err := binary.Write(&indexBuf, binary.BigEndian, uint16(len(labelBytes))); err != nil {
			return err
		}
		if _, err := indexBuf.Write(labelBytes); err != nil {
			return err
		}
	}

	n, err := crypto.EncryptStream(w, &indexBuf, keys)
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
	Truncate(int64) error
}

func (v *Vault) Save(w Writer) error {
	if !v.Full {
		for _, s := range v.orphanedSpans {
			offset := int64(v.HeaderLength() + s.Start)
			if _, err := w.Seek(offset, io.SeekStart); err != nil {
				return err
			}
			zeroBuf := make([]byte, s.Length)
			if _, err := w.Write(zeroBuf); err != nil {
				return err
			}
		}
	}
	v.orphanedSpans = nil

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

	pos, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if err := w.Truncate(pos); err != nil {
		return err
	}

	return nil
}

func (v *Vault) saveHeader(w io.Writer) error {
	if _, err := w.Write([]byte("VLT1")); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(len(v.devices))); err != nil {
		return err
	}
	for _, dev := range v.devices {
		if _, err := w.Write(dev.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
