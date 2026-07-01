//go:build !unix && !windows

package mem

func Lock(buf []byte) error  { return nil }
func Unlock(buf []byte) error { return nil }
