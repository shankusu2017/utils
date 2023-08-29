package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func getIV16() []byte {
	iv, _ := hex.DecodeString("daca11ed1f3fc59b2d233ec67cc6f028")
	return iv[:]
}

func getKey16() []byte {
	key, _ := hex.DecodeString("ac5299e1424c188fdb618ee0ee5481f7")
	return key[:]
}

func AESCrypt(b, iv, key []byte) []byte {
	plaintext := Pkcs7(b[:], aes.BlockSize)

	// encrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv[:])
	cipherText := make([]byte, len(plaintext))
	mode.CryptBlocks(cipherText[:], plaintext)

	return cipherText
}

func AESDeCrypt(cipherText, iv, key []byte) ([]byte, error) {
	if len(cipherText)%aes.BlockSize != 0 {
		err := errors.New(fmt.Sprintf("9eb49376 invalid ASE.buf.len:%d\n", len(cipherText)))
		return nil, err
	}

	// encrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// decrypt
	mode := cipher.NewCBCDecrypter(block, iv[:])
	plaintext := make([]byte, len(cipherText))
	mode.CryptBlocks(plaintext, cipherText[:])

	finalText := DePkcs7(plaintext, aes.BlockSize)

	return finalText, nil
}
