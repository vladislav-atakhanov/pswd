package crypto

import (
	"crypto/ecdh"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"golang.org/x/crypto/chacha20poly1305"
)

// DecryptStream decrypts an encrypted stream for a specific device.
// totalDevices is validated against the stream's device count;
// myDeviceIndex selects which key slot to decrypt (0, 1, 2...).
func DecryptStream(out io.Writer, encryptedStream io.Reader, privBytes [32]byte, totalDevices int, myDeviceIndex int) error {
	defer mem.ZeroArray32(&privBytes)

	ephemeralPubBytes := make([]byte, 32)
	if _, err := io.ReadFull(encryptedStream, ephemeralPubBytes); err != nil {
		return fmt.Errorf("read ephemeral key: %w", err)
	}

	ephemeralPub, err := ecdh.X25519().NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return fmt.Errorf("invalid ephemeral key: %w", err)
	}

	myPriv, err := ecdh.X25519().NewPrivateKey(privBytes[:])
	if err != nil {
		return err
	}
	sharedSecret, err := myPriv.ECDH(ephemeralPub)
	if err != nil {
		return fmt.Errorf("ECDH: %w", err)
	}
	mem.Lock(sharedSecret)
	defer mem.Unlock(sharedSecret)
	defer mem.ZeroBytes(sharedSecret)

	var numDevices uint16
	if err := binary.Read(encryptedStream, binary.BigEndian, &numDevices); err != nil {
		return fmt.Errorf("read device count: %w", err)
	}
	if int(numDevices) != totalDevices {
		return fmt.Errorf("device count mismatch: stream has %d devices, caller specified %d", numDevices, totalDevices)
	}

	slotSize := chacha20poly1305.NonceSize + 32 + chacha20poly1305.Overhead
	slotsBuf := make([]byte, int(numDevices)*slotSize)
	if _, err := io.ReadFull(encryptedStream, slotsBuf); err != nil {
		return fmt.Errorf("read key slots: %w", err)
	}

	if myDeviceIndex < 0 || myDeviceIndex >= int(numDevices) {
		return fmt.Errorf("device index %d out of range (0..%d)", myDeviceIndex, numDevices-1)
	}
	mySlot := slotsBuf[myDeviceIndex*slotSize : (myDeviceIndex+1)*slotSize]
	myNonce := mySlot[:chacha20poly1305.NonceSize]
	myCiphertext := mySlot[chacha20poly1305.NonceSize:]

	aead, err := chacha20poly1305.New(sharedSecret)
	if err != nil {
		return err
	}
	dataKey, err := aead.Open(nil, myNonce, myCiphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt data key: %w", err)
	}
	mem.Lock(dataKey)
	defer mem.Unlock(dataKey)
	defer mem.ZeroBytes(dataKey)

	mainAead, err := chacha20poly1305.New(dataKey)
	if err != nil {
		return err
	}

	mainNonce := make([]byte, mainAead.NonceSize())
	if _, err := io.ReadFull(encryptedStream, mainNonce); err != nil {
		return fmt.Errorf("read main nonce: %w", err)
	}

	ciphertext, err := io.ReadAll(encryptedStream)
	if err != nil {
		return fmt.Errorf("read ciphertext: %w", err)
	}

	plaintext, err := mainAead.Open(nil, mainNonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt content: %w", err)
	}

	_, err = out.Write(plaintext)
	return err
}
