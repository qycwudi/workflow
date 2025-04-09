package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
)

func Md5(str string) string {
	h := md5.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(str string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(decoded)
}

// 16, 24, 或 32
func GenAesKey(length int) string {
	key := make([]byte, length)
	_, _ = rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func AesEncrypt(str, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	// 创建密文缓冲区，包含 IV 和实际密文
	ciphertext := make([]byte, aes.BlockSize+len(str))
	iv := ciphertext[:aes.BlockSize]

	// 生成随机 IV
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	// 创建 CFB 加密器
	cfb := cipher.NewCFBEncrypter(block, iv)

	// 加密数据
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(str))

	// 返回 base64 编码的完整密文（包含 IV）
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func AesDecrypt(str, key string) string {
	// 解码 base64 字符串
	ciphertext, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	// 检查密文长度是否足够包含 IV
	if len(ciphertext) < aes.BlockSize {
		return ""
	}

	// 提取 IV 和实际密文
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建 CFB 解密器
	cfb := cipher.NewCFBDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)

	return string(plaintext)
}

// GenRsaKey 生成 RSA 密钥对，返回 PEM 格式的私钥和公钥
func GenRsaKey() map[string]string {
	// 生成 2048 位的 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return map[string]string{}
	}

	// 将私钥转换为 PKCS#1 格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// 将私钥编码为 PEM 格式
	privateKeyPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))

	// 将公钥转换为 PKIX 格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return map[string]string{}
	}

	// 将公钥编码为 PEM 格式
	publicKeyPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))
	return map[string]string{
		"privateKey": privateKeyPEM,
		"publicKey":  publicKeyPEM,
	}
}

// RsaEncrypt 使用 RSA 公钥加密数据
func RsaEncrypt(plaintext, publicKeyPEM string) string {
	// 解码 PEM 格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "failed to decode public key"
	}

	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err.Error()
	}

	// 类型断言为 RSA 公钥
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "failed to assert public key"
	}

	// 加密数据
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(plaintext))
	if err != nil {
		return err.Error()
	}

	// 返回 base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// RsaDecrypt 使用 RSA 私钥解密数据
func RsaDecrypt(ciphertext, privateKeyPEM string) string {
	// 解码 base64 密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return err.Error()
	}

	// 解码 PEM 格式的私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "failed to decode private key"
	}

	// 解析私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err.Error()
	}

	// 解密数据
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertextBytes)
	if err != nil {
		return err.Error()
	}

	return string(plaintext)
}

// PemToBase64 将 PEM 格式的密钥转换为 base64 字符串（用于环境变量）
func PemToBase64(pemStr string) string {
	return base64.StdEncoding.EncodeToString([]byte(pemStr))
}

// Base64ToPem 将 base64 字符串转换回 PEM 格式
func Base64ToPem(base64Str string) string {
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "failed to decode base64 string"
	}
	return string(bytes)
}

// PemToPrivateKey 将 PEM 格式的私钥字符串转换为 *rsa.PrivateKey
func PemToPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	// 解码 PEM 格式的私钥
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// SignData 使用 RSA 私钥对数据进行签名
func SignData(data string, privateKey *rsa.PrivateKey) (string, error) {
	// 计算数据的 SHA256 哈希值
	hash := sha256.Sum256([]byte(data))

	// 使用私钥对哈希值进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	// 返回 base64 编码的签名
	return base64.StdEncoding.EncodeToString(signature), nil
}
