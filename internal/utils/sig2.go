package utils

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	// Triple DES key size (in bytes)
	tripleDESKeySize = 24
	// Triple DES block size (in bytes)
	tripleDESBlockSize = 8
)

// pkcs7Pad pads the data to the block size using PKCS#7 padding.
func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, errors.New("invalid block size")
	}
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...), nil
}

// pkcs7Unpad removes PKCS#7 padding from the data.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, errors.New("invalid block size")
	}
	length := len(data)
	if length == 0 {
		return nil, errors.New("cannot unpad empty data")
	}
	if length%blockSize != 0 {
		// Data length must be a multiple of the block size for CBC decryption + padding
		// However, some implementations might allow this if the original data wasn't padded.
		// For strict PKCS#7 compatibility after CBC, it should be a multiple.
		// Let's return an error for clarity, matching typical expectations.
		return nil, errors.New("decrypted data length is not a multiple of block size")
	}

	padding := int(data[length-1])
	// Padding value cannot be zero or greater than block size
	if padding == 0 || padding > blockSize {
		// This might indicate bad padding or data that wasn't padded
		// For PKCS#7, this is invalid.
		return nil, errors.New("invalid padding value")
	}

	// Check that all padding bytes have the correct value
	padBytes := data[length-padding:]
	for _, b := range padBytes {
		if int(b) != padding {
			return nil, errors.New("invalid padding bytes")
		}
	}

	return data[:length-padding], nil
}

// Encrypt3DES performs 3DES encryption using CBC mode and PKCS#7 padding,
// then encodes the result in Base64.
// key must be 24 bytes.
// vector (IV) must be 8 bytes.
func Encrypt3DES(toEncode string, key string, vector string) string {
	keyBytes := []byte(key)
	ivBytes := []byte(vector)

	// --- Validation ---
	if len(keyBytes) != tripleDESKeySize {
		return fmt.Errorf("invalid key size: expected %d bytes, got %d", tripleDESKeySize, len(keyBytes)).Error()
	}
	if len(ivBytes) != tripleDESBlockSize {
		return fmt.Errorf("invalid IV size: expected %d bytes, got %d", tripleDESBlockSize, len(ivBytes)).Error()
	}
	// --- End Validation ---

	plaintextBytes := []byte(toEncode) // Assumes UTF-8 input string

	// Create 3DES cipher block
	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to create 3DES cipher: %w", err).Error()
	}

	// Check block size (should always be 8 for 3DES)
	blockSize := block.BlockSize()
	if blockSize != tripleDESBlockSize {
		// This should technically never happen for des.NewTripleDESCipher
		return fmt.Errorf("unexpected 3DES block size: %d", blockSize).Error()
	}

	// Apply PKCS#7 padding
	paddedPlaintext, err := pkcs7Pad(plaintextBytes, blockSize)
	if err != nil {
		return fmt.Errorf("failed to pad plaintext: %w", err).Error()
	}

	// Create CBC encrypter
	mode := cipher.NewCBCEncrypter(block, ivBytes)

	// Encrypt the data
	encryptedBytes := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(encryptedBytes, paddedPlaintext)

	// Encode result to Base64
	encodedString := base64.StdEncoding.EncodeToString(encryptedBytes)

	return encodedString
}

// Decrypt3DES decodes Base64 input, then performs 3DES decryption
// using CBC mode and removes PKCS#7 padding.
// key must be 24 bytes.
// vector (IV) must be 8 bytes.
func Decrypt3DES(toDecode string, key string, vector string) string {
	keyBytes := []byte(key)
	ivBytes := []byte(vector)

	// --- Validation ---
	if len(keyBytes) != tripleDESKeySize {
		return fmt.Errorf("invalid key size: expected %d bytes, got %d", tripleDESKeySize, len(keyBytes)).Error()
	}
	if len(ivBytes) != tripleDESBlockSize {
		return fmt.Errorf("invalid IV size: expected %d bytes, got %d", tripleDESBlockSize, len(ivBytes)).Error()
	}
	// --- End Validation ---

	// Decode Base64 input
	ciphertextBytes, err := base64.StdEncoding.DecodeString(toDecode)
	if err != nil {
		return fmt.Errorf("failed to decode base64 input: %w", err).Error()
	}

	// Create 3DES cipher block
	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to create 3DES cipher: %w", err).Error()
	}

	// Check block size
	blockSize := block.BlockSize()
	if blockSize != tripleDESBlockSize {
		return fmt.Errorf("unexpected 3DES block size: %d", blockSize).Error()
	}

	// Check if ciphertext length is valid (must be multiple of block size for CBC)
	if len(ciphertextBytes)%blockSize != 0 {
		return fmt.Errorf("ciphertext length is not a multiple of the block size").Error()
	}

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	// Decrypt the data
	// We can decrypt in-place or create a new buffer. Let's create a new one.
	decryptedPaddedBytes := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(decryptedPaddedBytes, ciphertextBytes)

	// Remove PKCS#7 padding
	decryptedBytes, err := pkcs7Unpad(decryptedPaddedBytes, blockSize)
	if err != nil {
		return fmt.Errorf("failed to unpad decrypted data: %w", err).Error()
	}

	// Convert result bytes to UTF-8 string
	decodedString := string(decryptedBytes)

	return decodedString
}

func DesSignData(datas string, privateKeyInput string) string {

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyInput)
	if err != nil {
		return fmt.Errorf("failed to decode base64 private key: %w", err).Error()
	}

	key, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse PKCS#8 private key: %w", err).Error()
	}

	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not an RSA key: %w", err).Error()
	}

	dataBytes := []byte(datas)

	hasher := sha256.New()
	_, err = hasher.Write(dataBytes)
	if err != nil {
		return fmt.Errorf("failed to hash data: %w", err).Error()
	}
	hashed := hasher.Sum(nil)

	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return fmt.Errorf("failed to sign data: %w", err).Error()
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	return signatureBase64
}

func SignRSA_MD5(data string, privateKeyB64 string) string {
	// 1. Decode Base64 Private Key
	keyBytes, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 private key: %w", err).Error()
	}

	// 2. Parse PKCS#8 Private Key
	// Use x509.ParsePKCS8PrivateKey for PKCS#8 format.
	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		// Optional: Try parsing as PKCS#1 if PKCS#8 fails, though the Java code specifies PKCS#8
		// parsedKey, err = x509.ParsePKCS1PrivateKey(keyBytes)
		// if err != nil {
		return fmt.Errorf("failed to parse PKCS#8 private key: %w", err).Error()
		// }
	}

	// Type assert to ensure it's an RSA private key
	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("key is not an RSA private key").Error()
	}

	// 3. Hash the data using MD5
	hash := md5.New()
	_, err = hash.Write([]byte(data)) // Use io.Writer interface for hashing
	if err != nil {
		// This specific error is unlikely with md5.New() writing bytes
		return fmt.Errorf("failed to write data to hash: %w", err).Error()
	}
	hashed := hash.Sum(nil) // Get the resulting hash bytes

	// 4. Sign the hash using RSA PKCS#1 v1.5
	// crypto.MD5 identifies the hash algorithm used.
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.MD5, hashed)
	if err != nil {
		return fmt.Errorf("failed to sign data: %w", err).Error()
	}

	// 5. Encode the signature to Base64
	signatureB64 := base64.StdEncoding.EncodeToString(signatureBytes)
	return signatureB64
}
