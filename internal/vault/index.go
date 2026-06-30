package vault

import (
	"encoding/binary"
	"io"
)

type IndexEntry struct{ raw []byte }

func (i IndexEntry) Bytes() []byte {
	return i.raw
}

const (
	uuidSize     = 16
	uint64Size   = 8
	labelLenSize = 2

	uuidOffset      = uuidSize
	updatedAtOffset = uuidSize
	startOffset     = updatedAtOffset + uint64Size
	lengthOffset    = startOffset + uint64Size
)

func (i *IndexEntry) UUID() [16]byte {
	var res [16]byte
	copy(res[:], i.raw[len(i.raw)-uuidOffset:])
	return res
}
func (i *IndexEntry) uint64(offset int) uint64 {
	end := len(i.raw) - offset
	return binary.LittleEndian.Uint64(i.raw[end-8 : end])
}
func (i *IndexEntry) UpdatedAt() uint64 {
	return i.uint64(updatedAtOffset)
}
func (i *IndexEntry) Start() int {
	return int(i.uint64(startOffset))
}
func (i *IndexEntry) Length() int {
	return int(i.uint64(lengthOffset))
}
func (i *IndexEntry) Label() string {
	labelLen := int(binary.BigEndian.Uint16(
		i.raw[len(i.raw)-42 : len(i.raw)-40],
	))
	return string(i.raw[len(i.raw)-42-labelLen : len(i.raw)-42])
}
func readIndex(r io.ReadSeeker) (IndexEntry, error) {
	const footerSize = 42 // 2 + 8 + 8 + 8 + 16

	end, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return IndexEntry{}, err
	}

	// Читаем хвост записи
	if _, err := r.Seek(end-footerSize, io.SeekStart); err != nil {
		return IndexEntry{}, err
	}

	var footer [footerSize]byte
	if _, err := io.ReadFull(r, footer[:]); err != nil {
		return IndexEntry{}, err
	}

	labelLen := int(binary.BigEndian.Uint16(footer[:2]))
	entrySize := footerSize + labelLen

	// Читаем запись целиком
	if _, err := r.Seek(end-int64(entrySize), io.SeekStart); err != nil {
		return IndexEntry{}, err
	}

	raw := make([]byte, entrySize)
	if _, err := io.ReadFull(r, raw); err != nil {
		return IndexEntry{}, err
	}
	return IndexEntry{raw: raw}, nil
}
func readIndexes(r io.ReadSeeker) ([]IndexEntry, error) {
	if _, err := r.Seek(-4, io.SeekEnd); err != nil {
		return nil, err
	}
	var count uint32
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return nil, err
	}
	if _, err := r.Seek(-4, io.SeekEnd); err != nil {
		return nil, err
	}
	res := make([]IndexEntry, count)
	for i := range count {
		index, err := readIndex(r)
		if err != nil {
			return nil, err
		}
		res[int(i)] = index
		if _, err := r.Seek(-int64(len(index.raw)), io.SeekCurrent); err != nil {
			return nil, err
		}
	}
	return res, nil
}
func newIndex(
	label string,
	length int,
	start int,
	updatedAt uint64,
	uuid [16]byte,
) IndexEntry {
	labelBytes := []byte(label)
	raw := make([]byte, len(labelBytes)+2+8+8+8+16)
	pos := 0
	copy(raw[pos:], labelBytes)
	pos += len(labelBytes)
	binary.BigEndian.PutUint16(raw[pos:], uint16(len(labelBytes)))
	pos += 2
	binary.LittleEndian.PutUint64(raw[pos:], uint64(length))
	pos += 8
	binary.LittleEndian.PutUint64(raw[pos:], uint64(start))
	pos += 8
	binary.LittleEndian.PutUint64(raw[pos:], updatedAt)
	pos += 8
	copy(raw[pos:], uuid[:])
	return IndexEntry{raw: raw}
}
