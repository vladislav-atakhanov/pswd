package crypto

import (
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/chacha20"
)

// EncryptStream принимает plainText как поток и пишет зашифрованный результат в out
func EncryptStream(out io.Writer, plainText io.Reader, publicKeys [][32]byte) error {
	// 1. Генерируем случайный 32-байтный ключ данных (Data Key)
	dataKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
		return err
	}

	// 2. Генерируем эфемерный ключ X25519 для ECDH
	ephemeralPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	ephemeralPubBytes := ephemeralPriv.PublicKey().Bytes()

	// Записываем эфемерный публичный ключ в самое начало потока
	if _, err := out.Write(ephemeralPubBytes); err != nil {
		return err
	}

	// Инициализируем HMAC для подписи всего, что пойдет в out (кроме самого HMAC, естественно)
	// В качестве ключа HMAC используем хэш от нашего dataKey, чтобы разграничить роли ключей
	macKey := sha256.Sum256(append(dataKey, []byte("mac")...))
	hmacSigner := hmac.New(sha256.New, macKey[:])

	// Оборачиваем out в MultiWriter, чтобы всё, что пишется на диск/в сеть, автоматически шло в HMAC
	teeOut := io.MultiWriter(out, hmacSigner)

	// 3. Шифруем dataKey для каждого устройства и пишем в поток
	for _, pubBytes := range publicKeys {
		devicePub, err := ecdh.X25519().NewPublicKey(pubBytes[:])
		if err != nil {
			return err
		}

		sharedSecret, err := ephemeralPriv.ECDH(devicePub)
		if err != nil {
			return err
		}

		// Шифруем dataKey с помощью XOR-потока ChaCha20 на базе sharedSecret
		keyNonce := make([]byte, chacha20.NonceSize) // Для эфемерного ключа nonce может быть нулевым
		cipherKey, err := chacha20.NewUnauthenticatedCipher(sharedSecret, keyNonce)
		if err != nil {
			return err
		}

		encryptedKey := make([]byte, 32)
		cipherKey.XORKeyStream(encryptedKey, dataKey)

		// Пишем зашифрованный ключ устройства в teeOut (уходит на диск и в HMAC)
		if _, err := teeOut.Write(encryptedKey); err != nil {
			return err
		}
	}

	// 4. Генерируем Nonce для шифрования основного тела данных
	mainNonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(rand.Reader, mainNonce); err != nil {
		return err
	}
	if _, err := teeOut.Write(mainNonce); err != nil {
		return err
	}

	// Инициализируем потоковый шифр для тела
	mainCipher, err := chacha20.NewUnauthenticatedCipher(dataKey, mainNonce)
	if err != nil {
		return err
	}

	// 5. Потоковое чтение из plainText, шифрование и запись в teeOut
	// Используем небольшой буфер (например, 4КБ), идеальный для RAM ограничений
	buf := make([]byte, 4096)
	for {
		n, err := plainText.Read(buf)
		if n > 0 {
			// Шифруем буфер на месте (in-place)
			mainCipher.XORKeyStream(buf[:n], buf[:n])
			// Пишем зашифрованный кусок в поток
			if _, err := teeOut.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// 6. Финальный штрих: вычисляем HMAC-подпись и дописываем её в самый конец
	// Она подписывает всё: зашифрованные ключи, nonce и шифротекст тела
	signature := hmacSigner.Sum(nil)
	if _, err := out.Write(signature); err != nil {
		return err
	}

	return nil
}
