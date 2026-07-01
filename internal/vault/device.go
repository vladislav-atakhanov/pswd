package vault

import (
	"bytes"
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

func (v *Vault) AddDevice(newPublicKey [32]byte, label string, r io.ReaderAt, privateKey [32]byte) error {
	v.Full = true
	for _, d := range v.Devices {
		if d.Name() == label {
			return fmt.Errorf("%w: %s", ErrDeviceExists, label)
		}
	}
	content := make(map[contentKey]Item, len(v.content))
	for id, c := range v.content {
		if c.content != nil {
			content[id] = c
			continue
		}
		b, err := v.Read(r, id, privateKey)
		if err != nil {
			return err
		}
		pass, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		content[id] = Item{
			Label:      c.Label,
			LastUpdate: c.LastUpdate,
			content:    bytes.NewReader(pass),
		}
	}
	v.content = content
	v.Devices = append(v.Devices, newDevice(newPublicKey, label))
	return nil
}

func (v *Vault) RemoveDevice(publicKey [32]byte, r io.ReaderAt, privateKey [32]byte) error {
	index := -1
	for i, d := range v.Devices {
		p := d.PublicKey()
		if bytes.Equal(p[:], publicKey[:]) {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrDeviceNotFound
	}

	v.Full = true

	content := make(map[contentKey]Item, len(v.content))
	for id, c := range v.content {
		if c.content != nil {
			content[id] = c
			continue
		}
		b, err := v.Read(r, id, privateKey)
		if err != nil {
			return err
		}
		pass, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		content[id] = Item{
			Label:      c.Label,
			LastUpdate: c.LastUpdate,
			content:    bytes.NewReader(pass),
		}
	}
	v.content = content

	v.Devices = append(v.Devices[:index], v.Devices[index+1:]...)
	return nil
}
