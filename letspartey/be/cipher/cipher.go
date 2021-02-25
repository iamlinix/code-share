package cipher

import (
	"encoding/base64"
	"errors"

	"iamlinix.com/partay/logger"
)

var _key string

func SetGlobalKey(key string) {
	_key = key
}

func DefaultEncrypt(msg string) (string, error) {
	return Encrypt(msg, _key)
}

func DefaultDecrypt(crypted string) (string, error) {
	return Decrypt(crypted, _key)
}

func CutInHalf(original []byte) []byte {
	len := len(original)
	if len%2 != 0 {
		original = append(original, 137)
		len += 1
	}
	halfLen := len / 2
	preHalf := original[:halfLen]
	postHalf := original[halfLen:]
	for i, v := range preHalf {
		postHalf[i] += v
	}
	return postHalf
}

func BendTheByte(b byte) byte {
	return (b + 77) % 121
}

func Encrypt(msg, key string) (encrypted string, err error) {
	if len(msg) == 0 || len(key) == 0 {
		logger.Errorf("invalid key or msg to encrypt: [%s:%s]", msg, key)
		return "", errors.New("invalid key or msg")
	}

	msg64 := base64.URLEncoding.EncodeToString([]byte(msg))
	key64 := base64.URLEncoding.EncodeToString([]byte(key))
	msgBytes := []byte(msg64)
	keyBytes := []byte(key64)
	halfKey := CutInHalf(keyBytes)
	halfLen := len(halfKey)
	keyIndex := 0
	encryptedBytes := make([]byte, len(msgBytes))
	for i, b := range msgBytes {
		encryptedBytes[i] = b ^ BendTheByte(halfKey[keyIndex])
		keyIndex += 1
		if keyIndex >= halfLen {
			keyIndex = 0
		}
	}

	return base64.URLEncoding.EncodeToString(encryptedBytes), nil
}

func Decrypt(crypted, key string) (decrypted string, err error) {
	cryptedR64, err := base64.URLEncoding.DecodeString(crypted)
	if err != nil {
		logger.Errorf("failed to decode crypted string: %s", crypted)
		return "", errors.New("cannot decode string")
	}

	key64 := base64.URLEncoding.EncodeToString([]byte(key))
	keyBytes := []byte(key64)
	halfKey := CutInHalf(keyBytes)
	halfLen := len(halfKey)
	keyIndex := 0
	decryptedBytes := make([]byte, len(cryptedR64))
	for i, b := range cryptedR64 {
		decryptedBytes[i] = b ^ BendTheByte(halfKey[keyIndex])
		keyIndex += 1
		if keyIndex >= halfLen {
			keyIndex = 0
		}
	}

	decryptedBytesR64, err := base64.URLEncoding.DecodeString(string(decryptedBytes))
	if err != nil {
		logger.Errorf("failed to double decode message: %v", err)
		return "", err
	}

	return string(decryptedBytesR64), nil
}
