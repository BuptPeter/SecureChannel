package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//生成32位md5字串
func GetMd5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetSaltMD5(input, salt string) string {
	hash := md5.New()
	//salt = "salt123456" //盐值
	io.WriteString(hash, input+salt)
	result := fmt.Sprintf("%x", hash.Sum(nil))
	return result
}

func GetCipherText(input []byte, key string) []byte {
	var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(input))
	cfb.XORKeyStream(ciphertext, input)
	return ciphertext
}

func GetPlainText(input []byte, key string) []byte {
	var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(input))
	cfbdec.XORKeyStream(plaintextCopy, input)
	return plaintextCopy
}
