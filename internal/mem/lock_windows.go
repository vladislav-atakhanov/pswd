//go:build windows

package mem

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func Lock(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return windows.VirtualLock(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}

func Unlock(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return windows.VirtualUnlock(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}
