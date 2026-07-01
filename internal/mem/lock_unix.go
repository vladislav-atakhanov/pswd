//go:build unix

package mem

import "golang.org/x/sys/unix"

func Lock(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return unix.Mlock(buf)
}

func Unlock(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return unix.Munlock(buf)
}
