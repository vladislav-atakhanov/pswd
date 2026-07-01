package crypto

import (
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20"
)

// DecryptStream принимает зашифрованный поток, расшифровывает его для конкретного девайса и пишет результат в out.
// totalDevices — общее количество устройств, записанное в заголовке файла.
// myDeviceIndex — индекс текущего устройства в списке (0, 1, 2...), чтобы понять, какой ключ забирать.
func DecryptStream(out io.Writer, encryptedStream io.Reader, privBytes [32]byte, totalDevices int, myDeviceIndex int) error {
	// 1. Читаем эфемерный публичный ключ (первые 32 байта)
	ephemeralPubBytes := make([]byte, 32)
	if _, err := io.ReadFull(encryptedStream, ephemeralPubBytes); err != nil {
		return fmt.Errorf("ошибка чтения эфемерного ключа: %w", err)
	}

	ephemeralPub, err := ecdh.X25519().NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return fmt.Errorf("неверный эфемерный ключ: %w", err)
	}

	// Восстанавливаем свой приватный ключ и считаем общий секрет (Shared Secret)
	myPriv, err := ecdh.X25519().NewPrivateKey(privBytes[:])
	if err != nil {
		return err
	}
	sharedSecret, err := myPriv.ECDH(ephemeralPub)
	if err != nil {
		return fmt.Errorf("ошибка ECDH: %w", err)
	}

	// 2. Инициализируем HMAC (пока без ключа данных, мы его узнаем через секунду)
	// Мы инициализируем его позже, но нам нужно логировать байты уже СЕЙЧАС.
	// Используем специальный io.Buffer для накопления байт, которые идут в HMAC,
	// либо обернем поток в io.TeeReader, как только расшифруем ключ.
	// Но проще собирать данные, идущие после эфемерного ключа, в HMAC вручную.

	var myEncryptedKey []byte

	// Нам нужно прочесть блоки ключей всех устройств (каждый по 32 байта)
	// И параллельно отправлять их в будущий HMAC. Для этого сохраним всю секцию ключей.
	allKeysLen := totalDevices * 32
	keysBuffer := make([]byte, allKeysLen)
	if _, err := io.ReadFull(encryptedStream, keysBuffer); err != nil {
		return fmt.Errorf("ошибка чтения ключей устройств: %w", err)
	}

	// Выдергиваем именно наш зашифрованный ключ данных
	myKeyOffset := myDeviceIndex * 32
	myEncryptedKey = keysBuffer[myKeyOffset : myKeyOffset+32]

	// 3. Расшифровываем основной Data Key
	keyNonce := make([]byte, chacha20.NonceSize)
	cipherKey, err := chacha20.NewUnauthenticatedCipher(sharedSecret, keyNonce)
	if err != nil {
		return err
	}

	dataKey := make([]byte, 32)
	cipherKey.XORKeyStream(dataKey, myEncryptedKey)

	// Теперь, когда у нас есть dataKey, мы можем запустить валидный HMAC!
	macKey := sha256.Sum256(append(dataKey, []byte("mac")...))
	hmacSigner := hmac.New(sha256.New, macKey[:])
	hmacSigner.Write(ephemeralPubBytes)

	// Пропускаем через HMAC уже прочитанную секцию ключей девайсов
	hmacSigner.Write(keysBuffer)

	// 4. Читаем Nonce для тела данных (12 байт) и отправляем в HMAC
	mainNonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(encryptedStream, mainNonce); err != nil {
		return fmt.Errorf("ошибка чтения nonce: %w", err)
	}
	hmacSigner.Write(mainNonce)

	// Инициализируем потоковый шифр для тела данных
	mainCipher, err := chacha20.NewUnauthenticatedCipher(dataKey, mainNonce)
	if err != nil {
		return err
	}

	// 5. Потоковое чтение тела с отсечением последних 32 байт (HMAC подписи)
	// Используем pending-буфер: накапливаем прочитанное, оставляя последние 32 байта нетронутыми.
	// Когда стрим заканчивается, в pending остаётся ровно HMAC-подпись (32 байта).
	var pending []byte
	chunk := make([]byte, 4096)
	for {
		n, err := encryptedStream.Read(chunk)
		if n > 0 {
			pending = append(pending, chunk[:n]...)
			for len(pending) > 32 {
				body := pending[:len(pending)-32]
				hmacSigner.Write(body)
				mainCipher.XORKeyStream(body, body)
				if _, err := out.Write(body); err != nil {
					return err
				}
				pending = pending[len(pending)-32:]
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// 6. ПРОВЕРКА ЦЕЛОСТНОСТИ: В pending сейчас лежит финальный HMAC с диска
	if len(pending) != 32 {
		return fmt.Errorf("поток слишком короткий: ожидалось 32 байта HMAC, получено %d", len(pending))
	}

	calculatedMac := hmacSigner.Sum(nil)

	if !hmac.Equal(pending, calculatedMac) {
		return errors.New("КРИТИЧЕСКАЯ ОШИБКА: Подпись HMAC не совпадает! Данные повреждены или подменены")
	}

	return nil
}
