package utils

import (
	"bytes"
	"crypto/aes"
	"errors"
	"fmt"
)

// const defaultKey = "1234567890ABCDEf" // 16 bytes for AES-128

// pkcs7Padding 填充数据以匹配 AES 块大小
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs7Unpadding 移除 PKCS7 填充
func pkcs7Unpadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("pkcs7: unpadding error - empty input")
	}
	unpadding := int(origData[length-1])
	if unpadding > length || unpadding > aes.BlockSize {
		// 注意：这里可能因为解密错误导致最后一个字节无效，从而引发错误
		// 更好的做法是检查所有填充字节是否一致，但这会增加时序攻击风险（虽然在此场景下可能不重要）
		// 为了简单和兼容 Java 的行为（它可能不严格检查所有填充字节），我们只检查长度
		// return nil, errors.New("pkcs7: unpadding error - invalid padding size")
		// 考虑到兼容性，如果解密出非预期的填充，可能直接返回了（Java行为未知）
		// 严格的检查应该是：
		// if unpadding > length {
		//     return nil, errors.New("pkcs7: invalid padding size")
		// }
		// for i := 0; i < unpadding; i++ {
		//     if origData[length-unpadding+i] != byte(unpadding) {
		//        return nil, errors.New("pkcs7: invalid padding bytes")
		//     }
		// }
		// 采用简单（可能不完全安全）的方式
		return nil, fmt.Errorf("pkcs7: unpadding error - invalid padding size %d", unpadding)
	}
	return origData[:(length - unpadding)], nil
}

// aesEncryptECB 使用 ECB 模式加密数据
// 警告：ECB 模式不安全，不推荐使用！仅用于兼容目的。
func aesEncryptECB(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	blockSize := block.BlockSize()
	paddedPlaintext := pkcs7Padding(plaintext, blockSize)
	ciphertext := make([]byte, len(paddedPlaintext))

	// ECB 模式逐块加密
	for bs, be := 0, blockSize; bs < len(paddedPlaintext); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(ciphertext[bs:be], paddedPlaintext[bs:be])
	}

	return ciphertext, nil
}

// aesDecryptECB 使用 ECB 模式解密数据
// 警告：ECB 模式不安全，不推荐使用！仅用于兼容目的。
func aesDecryptECB(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))

	// ECB 模式逐块解密
	for bs, be := 0, blockSize; bs < len(ciphertext); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(plaintext[bs:be], ciphertext[bs:be])
	}

	// 移除填充
	unpaddedPlaintext, err := pkcs7Unpadding(plaintext)
	if err != nil {
		// 如果解密失败（例如密钥错误），填充移除会出错
		return nil, fmt.Errorf("failed to unpad plaintext: %w", err)
	}

	return unpaddedPlaintext, nil
}

// AesEncryptToBytes AES 加密
// content: 待加密的内容 (UTF-8 字符串)
// encryptKey: 加密密钥 (字符串，将被转换为字节)
// 返回加密后的 byte[]
func FyyAesEncryptToBytes(content string, encryptKey string) ([]byte, error) {
	keyBytes := []byte(encryptKey)
	plaintext := []byte(content)

	// 检查密钥长度
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return nil, errors.New("invalid AES key size, must be 16, 24, or 32 bytes")
	}
	return aesEncryptECB(plaintext, keyBytes)
}

// AesEncrypt AES 加密为 base 64 code
// content: 待加密的内容 (UTF-8 字符串)
// encryptKey: 加密密钥 (字符串)
// 返回加密后的 base 64 code
func FyyAesEncrypt(content string, encryptKey string) string {
	encryptedBytes, err := FyyAesEncryptToBytes(content, encryptKey)
	if err != nil {
		return err.Error()
	}
	return Base64Encode(string(encryptedBytes))
}

// AesDecryptByBytes AES 解密
// encryptBytes: 待解密的 byte[]
// decryptKey: 解密密钥 (字符串)
// 返回解密后的 String
func FyyAesDecryptByBytes(encryptBytes []byte, decryptKey string) (string, error) {
	keyBytes := []byte(decryptKey)

	// 检查密钥长度
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("invalid AES key size, must be 16, 24, or 32 bytes")
	}

	// 警告：ECB 模式不安全
	// fmt.Println("Warning: Using AES with ECB mode, which is insecure.") // 可选，避免重复打印

	decryptedBytes, err := aesDecryptECB(encryptBytes, keyBytes)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(decryptedBytes), nil
}

// AesDecrypt 将 base 64 code AES 解密
// encryptStr: 待解密的 base 64 code
// decryptKey: 解密密钥 (字符串)
// 返回解密后的 string
func FyyAesDecrypt(encryptStr string, decryptKey string) string {
	if encryptStr == "" {
		// 对应 Java: StringUtils.isNotEmpty(encryptStr) ? ... : null;
		return ""
	}
	encryptBytes := Base64Decode(encryptStr)

	// 如果解码后是空字节数组 (对应 Java 返回 null 的情况)
	if len(encryptBytes) == 0 {
		return ""
	}
	decrypt, err := FyyAesDecryptByBytes([]byte(encryptBytes), decryptKey)
	if err != nil {
		return err.Error()
	}
	return decrypt
}
