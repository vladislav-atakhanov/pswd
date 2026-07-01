package uuid

import (
	"crypto/rand"
	"fmt"

	"encoding/hex"
	"errors"
)

type V4 [16]byte

func (u *V4) String() string {
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:16],
	)
}

func UUIDv4FromString(s string) (V4, error) {
	var u V4

	if len(s) != 36 {
		return u, errors.New("invalid UUID length")
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return u, errors.New("invalid UUID format")
	}

	if _, err := hex.Decode(u[0:4], []byte(s[0:8])); err != nil {
		return u, err
	}
	if _, err := hex.Decode(u[4:6], []byte(s[9:13])); err != nil {
		return u, err
	}
	if _, err := hex.Decode(u[6:8], []byte(s[14:18])); err != nil {
		return u, err
	}
	if _, err := hex.Decode(u[8:10], []byte(s[19:23])); err != nil {
		return u, err
	}
	if _, err := hex.Decode(u[10:16], []byte(s[24:36])); err != nil {
		return u, err
	}

	return u, nil
}

func NewV4() (V4, error) {
	var uuid V4
	if _, err := rand.Read(uuid[:]); err != nil {
		return uuid, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid, nil
}
