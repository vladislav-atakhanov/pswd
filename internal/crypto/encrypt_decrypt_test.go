package crypto

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func generateTestKeypair(t *testing.T) (priv, pub [32]byte) {
	t.Helper()
	priv, pub, err := GenerateKeys()
	if err != nil {
		t.Fatalf("GenerateKeys: %v", err)
	}
	return priv, pub
}

func encryptOut(t *testing.T, plain []byte, pubKeys [][32]byte) ([]byte, int) {
	t.Helper()
	var buf bytes.Buffer
	n, err := EncryptStream(&buf, bytes.NewReader(plain), pubKeys)
	if err != nil {
		t.Fatalf("EncryptStream: %v", err)
	}
	return buf.Bytes(), n
}

func decryptOut(t *testing.T, cipher []byte, priv [32]byte, totalDevices, idx int) ([]byte, error) {
	t.Helper()
	var out bytes.Buffer
	err := DecryptStream(&out, bytes.NewReader(cipher), priv, totalDevices, idx)
	return out.Bytes(), err
}

func TestRoundTrip(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte("Hello, World! This is a test message for round trip verification")

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatalf("round trip mismatch:\ngot:  %q\nwant: %q", result, plain)
	}
}

func TestMultipleDevices(t *testing.T) {
	n := 3
	privs := make([][32]byte, n)
	pubs := make([][32]byte, n)
	for i := range n {
		privs[i], pubs[i] = generateTestKeypair(t)
	}

	plain := []byte("Shared secret message for all devices")
	cipher, _ := encryptOut(t, plain, pubs)

	for i := range n {
		result, err := decryptOut(t, cipher, privs[i], n, i)
		if err != nil {
			t.Fatalf("device %d: DecryptStream: %v", i, err)
		}
		if !bytes.Equal(result, plain) {
			t.Fatalf("device %d: round trip mismatch:\ngot:  %q\nwant: %q", i, result, plain)
		}
	}
}

func TestEmptyPlaintext(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte{}

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatalf("round trip mismatch:\ngot:  %q\nwant: %q", result, plain)
	}
}

func TestLargeData(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := make([]byte, 10000)
	_, err := io.ReadFull(rand.Reader, plain)
	if err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatal("round trip mismatch for large data")
	}
}

func TestSmallPayload(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte("Hello!")

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatalf("round trip mismatch:\ngot:  %q\nwant: %q", result, plain)
	}
}

func TestEncryptStreamLength(t *testing.T) {
	_, pub := generateTestKeypair(t)
	plain := []byte("data")

	cipher, n := encryptOut(t, plain, [][32]byte{pub})
	if n != len(cipher) {
		t.Fatalf("returned length %d does not match written bytes %d", n, len(cipher))
	}
}

func TestExact32BytePayload(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte("This is exactly 32 bytes long!!!")

	if len(plain) != 32 {
		t.Fatalf("test data must be exactly 32 bytes, got %d", len(plain))
	}

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatalf("round trip mismatch:\ngot:  %q\nwant: %q", result, plain)
	}
}

func TestWrongPrivateKey(t *testing.T) {
	_, pubA := generateTestKeypair(t)
	privB, _ := generateTestKeypair(t)

	plain := []byte("This is a long enough message to avoid the small payload bug")
	cipher, _ := encryptOut(t, plain, [][32]byte{pubA})

	_, err := decryptOut(t, cipher, privB, 1, 0)
	if err == nil {
		t.Fatal("expected error for wrong private key, got nil")
	}
}

func TestTamperedCiphertext(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte("This is a long enough message to avoid the small payload bug")

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})

	cipher[len(cipher)-40] ^= 0x01

	_, err := decryptOut(t, cipher, priv, 1, 0)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext, got nil")
	}
}

func TestTruncatedStream(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := []byte("This is a long enough message to avoid the small payload bug")

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})

	truncated := cipher[:len(cipher)-10]

	var out bytes.Buffer
	err := DecryptStream(&out, bytes.NewReader(truncated), priv, 1, 0)
	if err == nil {
		t.Fatal("expected error for truncated stream, got nil")
	}
}

func TestWrongDeviceCount(t *testing.T) {
	priv0, pub0 := generateTestKeypair(t)
	_, pub1 := generateTestKeypair(t)

	plain := []byte("test")
	cipher, _ := encryptOut(t, plain, [][32]byte{pub0, pub1})

	var out bytes.Buffer
	err := DecryptStream(&out, bytes.NewReader(cipher), priv0, 3, 0)
	if err == nil {
		t.Fatal("expected error for wrong device count (3 instead of 2), got nil")
	}
}

func TestTooFewDevices(t *testing.T) {
	priv0, pub0 := generateTestKeypair(t)
	_, pub1 := generateTestKeypair(t)

	plain := []byte("test")
	cipher, _ := encryptOut(t, plain, [][32]byte{pub0, pub1})

	var out bytes.Buffer
	err := DecryptStream(&out, bytes.NewReader(cipher), priv0, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bytes.Equal(out.Bytes(), plain) {
		t.Fatal("expected decrypted output to differ from plaintext when totalDevices is wrong")
	}
}

func TestNonMultiplePayload(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := make([]byte, 5000)
	_, err := io.ReadFull(rand.Reader, plain)
	if err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatal("round trip mismatch for non-multiple payload")
	}
}

func TestPayloadBetween32And4096(t *testing.T) {
	sizes := []int{33, 100, 500, 1000, 2000, 4000}
	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			priv, pub := generateTestKeypair(t)
			plain := make([]byte, size)
			_, err := io.ReadFull(rand.Reader, plain)
			if err != nil {
				t.Fatalf("rand.Read: %v", err)
			}

			cipher, _ := encryptOut(t, plain, [][32]byte{pub})
			result, err := decryptOut(t, cipher, priv, 1, 0)
			if err != nil {
				t.Fatalf("DecryptStream: %v", err)
			}

			if !bytes.Equal(result, plain) {
				t.Fatal("round trip mismatch")
			}
		})
	}
}

func TestPayloadExact4096(t *testing.T) {
	priv, pub := generateTestKeypair(t)
	plain := make([]byte, 4096)
	_, err := io.ReadFull(rand.Reader, plain)
	if err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	cipher, _ := encryptOut(t, plain, [][32]byte{pub})
	result, err := decryptOut(t, cipher, priv, 1, 0)
	if err != nil {
		t.Fatalf("DecryptStream: %v", err)
	}

	if !bytes.Equal(result, plain) {
		t.Fatal("round trip mismatch for exact 4096-byte payload")
	}
}
