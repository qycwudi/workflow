package utils

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
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

// DESSign 3DES签名
func DESSign(datas, privates string) string {
	// Sign generates a signature for the given data using the provided private key.
	// datas: The string data to be signed.
	// privates: A Base64 encoded string of the PKCS#8 private key.
	// Returns the Base64 encoded signature string or an error.
	// 1. Decode the Base64 private key string
	// Assuming CipherBase64 uses standard Base64 encoding
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privates)
	if err != nil {
		return fmt.Errorf("failed to decode base64 private key: %w", err).Error()
	}

	// 2. Parse the PKCS#8 encoded private key
	// This step implicitly handles the KEY_ALGORITHM ("RSA" in our assumption)
	key, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		// Optional: If the key might be PKCS#1 instead of PKCS#8, you could try parsing it too.
		// key, err = x509.ParsePKCS1PrivateKey(privateKeyBytes)
		// if err != nil {
		//     return "", fmt.Errorf("failed to parse private key (tried PKCS8 and PKCS1): %w", err)
		// }
		return fmt.Errorf("failed to parse PKCS#8 private key: %w", err).Error()
	}

	// Type assert the key to the specific type (*rsa.PrivateKey in our assumption)
	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		// If KEY_ALGORITHM could be something else (like EC), you'd need to handle other types:
		// ecdsaPrivateKey, ok := key.(*ecdsa.PrivateKey)
		// if !ok { ... }
		return fmt.Errorf("private key is not an RSA key: %w", err).Error()
	}

	// 3. Prepare the data to be signed
	dataBytes := []byte(datas)

	// 4. Hash the data
	// The SIGNATURE_ALGORITHM ("SHA256withRSA") determines the hash function.
	hasher := sha256.New() // Use sha1.New() if the algorithm was SHA1withRSA
	_, err = hasher.Write(dataBytes)
	if err != nil {
		return fmt.Errorf("failed to hash data: %w", err).Error()
	}
	hashed := hasher.Sum(nil)

	// 5. Sign the hashed data using the private key
	// The crypto.SHA256 identifier corresponds to the hash algorithm used.
	// rsa.SignPKCS1v15 is commonly used for "SHAxxxwithRSA" signatures.
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return fmt.Errorf("failed to sign data: %w", err).Error()
	}

	// 6. Encode the signature bytes to Base64
	// Assuming CipherBase64 uses standard Base64 encoding
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	return signatureBase64

}
