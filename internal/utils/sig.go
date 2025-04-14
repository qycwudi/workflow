package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

// Algorithm constants
const (
	AlgorithmDES     = "desede"
	AlgorithmModeDES = "/CBC/PKCS5Padding"
)

// Encrypt 3DES加密
func DESEncrypt(toEncode, key, vector string) string {
	// 检查密钥长度
	if len(key) != 24 {
		return errors.New("key length must be 24 bytes").Error()
	}

	// 创建3DES cipher block
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return fmt.Errorf("failed to create cipher block: %v", err).Error()
	}

	// 检查IV长度
	if len(vector) < block.BlockSize() {
		return fmt.Errorf("vector length must be at least %d bytes", block.BlockSize()).Error()
	}
	iv := []byte(vector)[:block.BlockSize()]

	// 创建加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// PKCS5Padding填充
	paddedData := DESPKCS5Padding([]byte(toEncode), block.BlockSize())

	// 加密
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt 3DES解密
func DESDecrypt(toDecode, key, vector string) string {
	// 检查密钥长度
	if len(key) != 24 {
		return errors.New("key length must be 24 bytes").Error()
	}

	// Base64解码
	ciphertext, err := base64.StdEncoding.DecodeString(toDecode)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %v", err).Error()
	}

	// 创建3DES cipher block
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return fmt.Errorf("failed to create cipher block: %v", err).Error()
	}

	// 检查IV长度
	if len(vector) < block.BlockSize() {
		return fmt.Errorf("vector length must be at least %d bytes", block.BlockSize()).Error()
	}
	iv := []byte(vector)[:block.BlockSize()]

	// 创建解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS5Unpadding去除填充
	unpadded, err := DESPKCS5Unpadding(plaintext)
	if err != nil {
		return fmt.Errorf("failed to unpad data: %v", err).Error()
	}

	return string(unpadded)
}

// Generate3DesPrivateKey 生成3des私钥
func DESGenerate3DesPrivateKey() string {
	// 3DES密钥需要24字节
	key := make([]byte, 24)
	_, err := rand.Read(key)
	if err != nil {
		return fmt.Errorf("failed to generate random key: %v", err).Error()
	}

	// Base64编码并截取前24个字符
	encoded := base64.StdEncoding.EncodeToString(key)
	if len(encoded) < 24 {
		return errors.New("generated key is too short").Error()
	}
	return encoded[:24]
}

// PKCS5Padding PKCS5填充
func DESPKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS5Unpadding PKCS5去除填充
func DESPKCS5Unpadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}

	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}

	return src[:(length - unpadding)], nil
}
