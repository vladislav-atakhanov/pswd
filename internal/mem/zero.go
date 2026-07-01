package mem

import "runtime"

func ZeroBytes(s []byte) {
	for i := range s {
		s[i] = 0
	}
	runtime.KeepAlive(s)
}

func ZeroArray32(a *[32]byte) {
	if a == nil {
		return
	}
	clear(a[:])
	runtime.KeepAlive(a)
}
