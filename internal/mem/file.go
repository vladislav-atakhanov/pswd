package mem

import (
	"errors"
	"io"
)

type MemoryFile struct {
	data []byte
	pos  int64
}

func (f *MemoryFile) Read(p []byte) (int, error) {
	if f.pos >= int64(len(f.data)) {
		return 0, io.EOF
	}
	n := copy(p, f.data[f.pos:])
	f.pos += int64(n)
	return n, nil
}

func (f *MemoryFile) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(f.data)) {
		return 0, io.EOF
	}
	n := copy(p, f.data[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (f *MemoryFile) Write(p []byte) (int, error) {
	end := f.pos + int64(len(p))
	if end > int64(len(f.data)) {
		buf := make([]byte, end)
		copy(buf, f.data)
		f.data = buf
	}
	n := copy(f.data[f.pos:], p)
	f.pos += int64(n)
	return n, nil
}

func (f *MemoryFile) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = f.pos + offset
	case io.SeekEnd:
		abs = int64(len(f.data)) + offset
	default:
		return 0, errors.New("invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("negative position")
	}
	f.pos = abs
	return abs, nil
}

func (f *MemoryFile) Truncate(size int64) error {
	if size < 0 {
		return errors.New("negative size")
	}
	if size < int64(len(f.data)) {
		f.data = f.data[:size]
	} else if size > int64(len(f.data)) {
		buf := make([]byte, size)
		copy(buf, f.data)
		f.data = buf
	}
	if f.pos > size {
		f.pos = size
	}
	return nil
}

func (f *MemoryFile) Bytes() []byte {
	return f.data
}

func (f *MemoryFile) Len() int {
	return len(f.data)
}
