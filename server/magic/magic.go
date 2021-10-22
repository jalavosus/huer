package magic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"github.com/pkg/errors"
)

func Decrypt(magicKey, encryptedMagic string) (string, error) {
	key := func() []byte {
		k := sha256.Sum256([]byte(magicKey))
		return k[:]
	}()

	magicBytes, err := decrypt(key, []byte(encryptedMagic))
	if err != nil {
		return "", err
	}

	return string(magicBytes), nil
}

func decrypt(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create AES block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM")
	}

	if len(cipherText) < gcm.NonceSize() {
		return nil, errors.Errorf(
			"malformed cipherText; length %[1]d less than GCM nonce size %[2]d",
			len(cipherText),
			gcm.NonceSize(),
		)
	}

	return gcm.Open(
		nil,
		cipherText[:gcm.NonceSize()],
		cipherText[gcm.NonceSize():],
		nil,
	)
}