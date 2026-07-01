package mem

import "runtime"

func ZeroBytes(s []byte) {
	clear(s)
	runtime.KeepAlive(s)
}

func ZeroArray32(a *[32]byte) {
	if a == nil {
		return
	}
	clear(a[:])
	runtime.KeepAlive(a)
}
