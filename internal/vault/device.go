package vault

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Device struct{ raw []byte }

func newDevice(publicKey [32]byte, label string) Device {
	labelBytes := []byte(label)
	res := make([]byte, 32+2+len(labelBytes))
	copy(res[:32], publicKey[:])
	binary.BigEndian.PutUint16(res[32:34], uint16(len(labelBytes)))
	copy(res[34:], labelBytes)
	return Device{res}
}
func readDevice(r io.Reader) (Device, error) {
	header := make([]byte, 34)
	if _, err := io.ReadFull(r, header); err != nil {
		return Device{}, err
	}
	labelLen := binary.BigEndian.Uint16(header[32:34])
	res := make([]byte, 34+int(labelLen))
	copy(res, header)
	if _, err := io.ReadFull(r, res[34:]); err != nil {
		return Device{}, err
	}
	return Device{res}, nil
}
func (d Device) Bytes() []byte {
	return d.raw
}

func (d *Device) PublicKey() [32]byte {
	var dst [32]byte
	copy(dst[:], d.raw)
	return dst
}
func (d *Device) Name() string {
	labelLen := binary.BigEndian.Uint16(d.raw[32:34])
	return string(d.raw[34 : 34+int(labelLen)])
}

func readDevices(r io.ReadSeeker) ([]Device, error) {
	var count uint16
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return nil, err
	}
	res := make([]Device, count)
	for i := range count {
		device, err := readDevice(r)
		if err != nil {
			return nil, err
		}
		res[i] = device
	}
	return res, nil
}

func (v *Vault) AddDevice(publicKey [32]byte, label string) error {
	for _, d := range v.Devices {
		if d.Name() == label {
			return fmt.Errorf("Device %s already in vault", label)
		}
	}
	v.Devices = append(v.Devices, newDevice(publicKey, label))
	return nil
}
