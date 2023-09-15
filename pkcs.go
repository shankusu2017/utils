package utils

import (
	"bytes"
	"log"
)

func Pkcs7(plaintext []byte, blockSize int) []byte {
	plainLen := len(plaintext)
	paddingSize := blockSize - plainLen%blockSize

	finalText := make([]byte, plainLen+paddingSize)

	copy(finalText[:], plaintext)
	copy(finalText[plainLen:], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))

	return finalText
}

func DePkcs7(buf []byte, blockLen int) []byte {
	if len(buf) == 0 {
		log.Printf("ERROR 47f4e43b input buf is empty\n")
		return []byte{}
	}
	paddingSize := int(buf[len(buf)-1])
	if paddingSize > blockLen {
		return buf
	}

	plainSizes := len(buf) - paddingSize
	text := make([]byte, plainSizes)
	copy(text[:], buf[:plainSizes])
	return text
}
