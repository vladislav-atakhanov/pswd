package vault

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func (v *Vault) read(r io.ReadSeeker, size int, privateKey [32]byte) error {
	if _, err := r.Seek(-4, io.SeekEnd); err != nil {
		return err
	}

	var indexLen uint32
	if err := binary.Read(r, binary.BigEndian, &indexLen); err != nil {
		return err
	}
	if size < 4 {
		return nil
	}
	if indexLen == 0 || indexLen > uint32(size-4) || indexLen > uint32(size) {
		return fmt.Errorf("invalid index length: %d (file size: %d)", indexLen, size)
	}
	v.dataEnd = size - 4 - int(indexLen)

	if _, err := r.Seek(-int64(indexLen)-4, io.SeekEnd); err != nil {
		return err
	}

	encryptedIndex := make([]byte, indexLen)
	if _, err := io.ReadFull(r, encryptedIndex); err != nil {
		return err
	}

	var decryptedIndex bytes.Buffer
	if err := v.decrypt(bytes.NewReader(encryptedIndex), &decryptedIndex, privateKey); err != nil {
		return err
	}

	reader := bytes.NewReader(decryptedIndex.Bytes())

	var count uint32
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}

	v.content = make(map[contentKey]Item, count)
	for i := uint32(0); i < count; i++ {
		var id contentKey
		if _, err := io.ReadFull(reader, id[:]); err != nil {
			return err
		}

		var lastUpdate uint64
		if err := binary.Read(reader, binary.BigEndian, &lastUpdate); err != nil {
			return err
		}

		var start, length uint32
		if err := binary.Read(reader, binary.BigEndian, &start); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}

		var labelLen uint16
		if err := binary.Read(reader, binary.BigEndian, &labelLen); err != nil {
			return err
		}

		label := make([]byte, labelLen)
		if _, err := io.ReadFull(reader, label); err != nil {
			return err
		}

		v.content[id] = Item{
			Label:      string(label),
			LastUpdate: lastUpdate,
			start:      int(start),
			length:     int(length),
		}
	}
	return nil
}
