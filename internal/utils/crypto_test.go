package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"log"
	"testing"
)

func TestMd5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "21232f297a57a5a743894a0e4a801fc3" + "21232f297a57a5a74"}, want: "3c3d20cf4936b81600306b09ab1f6cf4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.str); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "admin"}, want: "YWRtaW4="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Encode(tt.args.str); got != tt.want {
				t.Errorf("Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestBase64Decode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "YWRtaW4="}, want: "admin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Decode(tt.args.str); got != tt.want {
				t.Errorf("Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesEncrypt(t *testing.T) {
	key := "QyJbR5bmiZSwhQjsXMivSA=="
	type args struct {
		str string
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "admin", key: key}, want: "cgBXsGCPnWF2DgidiU7CYIzV6nFU"},
	}
	fmt.Println(key)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AesEncrypt(tt.args.str, tt.args.key); got != tt.want {
				t.Errorf("AesEncrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesDecrypt(t *testing.T) {
	key := "QyJbR5bmiZSwhQjsXMivSA=="
	type args struct {
		str string
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "jZKTKAxPV1f4u+FICOjl8slOpQQs", key: key}, want: "admin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AesDecrypt(tt.args.str, tt.args.key); got != tt.want {
				t.Errorf("AesDecrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenRsaKey(t *testing.T) {
	keys := GenRsaKey()
	fmt.Println("------------私钥--------------------")
	fmt.Println(keys["privateKey"])
	fmt.Println("------------公钥--------------------")
	fmt.Println(keys["publicKey"])
}

/*

------------私钥--------------------
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEApnWyufKUqpxJ1BnrGBHY/qw2iZmwH7mQwknLKuy0GzwJE38b
15YyOnUQsixqnMjS+ijYoewnIjoJzkRG+IKILhwYfr/KICp+yUUOUdznj+UWcduE
4ArMJClmQO7QxH7i6aEUnZZaWR0Y6fNF8tNAAK6BMMXBweijr3DK8NbU4BicCSE6
cA0QFaa/AShZwbh0jHOfH7o21YKxp0OKi2m3/cT1eW49dxD4NICQWgHcyF5p5Q+I
ZGlWkAjDRX0tGWZn3I/DRrHLmOt5BkV/si/WUH14ohjkqxNaCspdagYh5N2W8y7e
QZ7OIVkHRkA1huShj1si2QW3ecf17AaGT26DFwIDAQABAoIBAFx8lrIkKx/kOE0F
nY7BE9zkWGT5pFbsJpccoyqVW7iqEazoedK92+iI9J7aeN3l4vslSMAm0tN8hcg6
PWiENh/QcneyHh7jHZC8sqkfYPWdZmwaeOh2d1g8PFkll48lpPQEEM1CJa7PsPED
vsCdSTlXu+Yaan5PwOqQc+F+q8ZaYFk7lmCEgOWVTgwA6M7DOp7rcEzTnuJY9cie
PfHOZRrPOT+/V01hiKsbUWYvb2MWbfF/wxJwRAsG4u7Z6sQUDYipcQ85CYxt6m+t
YOap2Az+MSnGeM8Pf+KA2EKoR9Q7JejcD056DVAJcMuQNDguWxz6nje0vWRlL79l
Haym5gkCgYEAw9AY88DomONfnLdYNNi12F5wQKKg4Uz46T+S0SK6KOrHs7J57VvW
4Kpy+jqHTDxpeCp2IOECBPT5nyMyfLsyBlFAVjw0b+i6STg1uw/8xyl8+Dma0v36
DfUwiCrCmRSTzlHWf3cKNtjheGNlFy8qIfWFBznIpnMZX3lgxDBhW4UCgYEA2Z/l
N8zIn7SNAVWcSCYuk2sy5LgkK4pozx9BXzcBlXai/EgD65DVYylQaEfl5QOTBE3K
et3iGGPJRAakW6pC6J9RsRO6HfGJdkuOeVSAKgm36Tm4rYfV0bu7hH+xssD0d9fd
2W8RynFtnkaaov9mnbIp43JHcVl0WnlVnFApgOsCgYEApeYYTeSB7I6vghJgXB3D
K3cPueNPVLMnLE8db6zxdhs8aQXsgWpPCne/BDw0RyXj4dhvzvl0AYkgOHDUpJLh
FjMexDEr6CiQM9q4wy0PaBnBdHkxsFNX2R2EKcm4p4OkmqgBiGrtr3xewuXLTzI5
ix39wBp34nYf6CDpGC85PRUCgYEArA40hSdMvqdai+GJi6lUTY0FUbscLahiMM7/
Oi4c/HQta9Pr9YQukRWK0sd1RNjMlSyDlxxxsuLBrxypOSelepDrX1q/XQknqvUV
kWtzYMkKNEREdD3emNEZ8ima7j6LiWyLo2qi4DFJf0dG3vOZx7eiUoZ5YW5eBWHE
g68FAT0CgYAgA4wOsUTVcNUm2eG903h8ogjhmPRXwnEg6qwmFXdhD3vWUzFx214m
Tr50O1kW46VxSXjTP/yVxolcSrlFzB15sZk5dHUaNiTj+KU4lNAV5IMc02JP2EbM
RY0jOBEaIAvHslZDNtyIccVKFLDjKtutc9BB0V/YvZvysoY9UxVJAA==
-----END RSA PRIVATE KEY-----

------------公钥--------------------
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApnWyufKUqpxJ1BnrGBHY
/qw2iZmwH7mQwknLKuy0GzwJE38b15YyOnUQsixqnMjS+ijYoewnIjoJzkRG+IKI
LhwYfr/KICp+yUUOUdznj+UWcduE4ArMJClmQO7QxH7i6aEUnZZaWR0Y6fNF8tNA
AK6BMMXBweijr3DK8NbU4BicCSE6cA0QFaa/AShZwbh0jHOfH7o21YKxp0OKi2m3
/cT1eW49dxD4NICQWgHcyF5p5Q+IZGlWkAjDRX0tGWZn3I/DRrHLmOt5BkV/si/W
UH14ohjkqxNaCspdagYh5N2W8y7eQZ7OIVkHRkA1huShj1si2QW3ecf17AaGT26D
FwIDAQAB
-----END PUBLIC KEY-----
*/

func TestRsaEncrypt(t *testing.T) {
	pub := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApnWyufKUqpxJ1BnrGBHY
/qw2iZmwH7mQwknLKuy0GzwJE38b15YyOnUQsixqnMjS+ijYoewnIjoJzkRG+IKI
LhwYfr/KICp+yUUOUdznj+UWcduE4ArMJClmQO7QxH7i6aEUnZZaWR0Y6fNF8tNA
AK6BMMXBweijr3DK8NbU4BicCSE6cA0QFaa/AShZwbh0jHOfH7o21YKxp0OKi2m3
/cT1eW49dxD4NICQWgHcyF5p5Q+IZGlWkAjDRX0tGWZn3I/DRrHLmOt5BkV/si/W
UH14ohjkqxNaCspdagYh5N2W8y7eQZ7OIVkHRkA1huShj1si2QW3ecf17AaGT26D
FwIDAQAB
-----END PUBLIC KEY-----`
	encrypt := RsaEncrypt("admin", pub)
	fmt.Println(encrypt)
	fmt.Println("--------------------------------")
	fmt.Println(PemToBase64(pub))
}

func TestRsaDecrypt(t *testing.T) {
	encrypt := "VcvewYRJMdaoigzwSiAUpIgfOVz5AoB6nkcTcCbYrsN7zi2VnB6rL3CbiTH6yQBFO0ZtiYCgq9pCLdRHig/jm1Q/6VG0DNfAob0jXx07/tkhElgXc3GHjQ4QbExCw6SPLdoOCBEhAJRoKqCRk4AizdDXtbstnroNaxYiuqIgvO/Iwr5j/WSpz3Kar6JL8C28pGm7YEsfTfc+xV512Wkyld2BnC1FFGaYkCvfMIqOSjE5yOY5Ljkpw8YsXnT7uQ9If7c6nFt1ggiQ1NbLH8kRAPHDjOAGlAg7XKupXtu/o+DS1Sv5Y0/Q121UdwaycDVepo1m0K82akwKkIyriHBRFg=="
	pri := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEApnWyufKUqpxJ1BnrGBHY/qw2iZmwH7mQwknLKuy0GzwJE38b
15YyOnUQsixqnMjS+ijYoewnIjoJzkRG+IKILhwYfr/KICp+yUUOUdznj+UWcduE
4ArMJClmQO7QxH7i6aEUnZZaWR0Y6fNF8tNAAK6BMMXBweijr3DK8NbU4BicCSE6
cA0QFaa/AShZwbh0jHOfH7o21YKxp0OKi2m3/cT1eW49dxD4NICQWgHcyF5p5Q+I
ZGlWkAjDRX0tGWZn3I/DRrHLmOt5BkV/si/WUH14ohjkqxNaCspdagYh5N2W8y7e
QZ7OIVkHRkA1huShj1si2QW3ecf17AaGT26DFwIDAQABAoIBAFx8lrIkKx/kOE0F
nY7BE9zkWGT5pFbsJpccoyqVW7iqEazoedK92+iI9J7aeN3l4vslSMAm0tN8hcg6
PWiENh/QcneyHh7jHZC8sqkfYPWdZmwaeOh2d1g8PFkll48lpPQEEM1CJa7PsPED
vsCdSTlXu+Yaan5PwOqQc+F+q8ZaYFk7lmCEgOWVTgwA6M7DOp7rcEzTnuJY9cie
PfHOZRrPOT+/V01hiKsbUWYvb2MWbfF/wxJwRAsG4u7Z6sQUDYipcQ85CYxt6m+t
YOap2Az+MSnGeM8Pf+KA2EKoR9Q7JejcD056DVAJcMuQNDguWxz6nje0vWRlL79l
Haym5gkCgYEAw9AY88DomONfnLdYNNi12F5wQKKg4Uz46T+S0SK6KOrHs7J57VvW
4Kpy+jqHTDxpeCp2IOECBPT5nyMyfLsyBlFAVjw0b+i6STg1uw/8xyl8+Dma0v36
DfUwiCrCmRSTzlHWf3cKNtjheGNlFy8qIfWFBznIpnMZX3lgxDBhW4UCgYEA2Z/l
N8zIn7SNAVWcSCYuk2sy5LgkK4pozx9BXzcBlXai/EgD65DVYylQaEfl5QOTBE3K
et3iGGPJRAakW6pC6J9RsRO6HfGJdkuOeVSAKgm36Tm4rYfV0bu7hH+xssD0d9fd
2W8RynFtnkaaov9mnbIp43JHcVl0WnlVnFApgOsCgYEApeYYTeSB7I6vghJgXB3D
K3cPueNPVLMnLE8db6zxdhs8aQXsgWpPCne/BDw0RyXj4dhvzvl0AYkgOHDUpJLh
FjMexDEr6CiQM9q4wy0PaBnBdHkxsFNX2R2EKcm4p4OkmqgBiGrtr3xewuXLTzI5
ix39wBp34nYf6CDpGC85PRUCgYEArA40hSdMvqdai+GJi6lUTY0FUbscLahiMM7/
Oi4c/HQta9Pr9YQukRWK0sd1RNjMlSyDlxxxsuLBrxypOSelepDrX1q/XQknqvUV
kWtzYMkKNEREdD3emNEZ8ima7j6LiWyLo2qi4DFJf0dG3vOZx7eiUoZ5YW5eBWHE
g68FAT0CgYAgA4wOsUTVcNUm2eG903h8ogjhmPRXwnEg6qwmFXdhD3vWUzFx214m
Tr50O1kW46VxSXjTP/yVxolcSrlFzB15sZk5dHUaNiTj+KU4lNAV5IMc02JP2EbM
RY0jOBEaIAvHslZDNtyIccVKFLDjKtutc9BB0V/YvZvysoY9UxVJAA==
-----END RSA PRIVATE KEY-----`
	decrypt := RsaDecrypt(encrypt, pri)
	fmt.Println(decrypt)
}

// Cipher3DES 加密函数
func Cipher3DESEncrypt(data, key, iv string) (string, error) {
	// 将字符串转换为字节数组
	plaintext := []byte(data)
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	// 检查密钥长度是否为 24 字节（3DES 需要 24 字节的密钥）
	if len(keyBytes) != 24 {
		return "", fmt.Errorf("invalid key length, expected 24 bytes")
	}

	// 检查 IV 长度是否为 8 字节
	if len(ivBytes) != 8 {
		return "", fmt.Errorf("invalid IV length, expected 8 bytes")
	}

	// 创建 3DES 加密块
	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create 3DES cipher block: %v", err)
	}

	// 填充明文数据（PKCS5 填充）
	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	// 创建 CBC 模式的加密器
	mode := cipher.NewCBCEncrypter(block, ivBytes)

	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 返回 Base64 编码的加密结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// PKCS5Padding 填充函数
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func TestCipher3DESEncrypt(t *testing.T) {
	// 原始数据
	data := `{"key": "value"}` // 替换为实际的 JSON 数据
	fmt.Println("请求报文：" + data)

	// 平台分配的唯一的接入秘钥
	AppKey := "l4mdofLTvHkyONpdlyXBiaTv"
	vector := "12345678" // 随机 8 位偏移量

	// 加密数据
	encrData, err := Cipher3DESEncrypt(data, AppKey, vector)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}
	fmt.Println("加密报文：" + encrData)
}

func TestSignData(t *testing.T) {
	// 假设 encrData 是上一步的加密结果
	encrData := "gzDzraJo5CisyMhC5dSCaWjae3yTIXW6"

	pri := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEApnWyufKUqpxJ1BnrGBHY/qw2iZmwH7mQwknLKuy0GzwJE38b
15YyOnUQsixqnMjS+ijYoewnIjoJzkRG+IKILhwYfr/KICp+yUUOUdznj+UWcduE
4ArMJClmQO7QxH7i6aEUnZZaWR0Y6fNF8tNAAK6BMMXBweijr3DK8NbU4BicCSE6
cA0QFaa/AShZwbh0jHOfH7o21YKxp0OKi2m3/cT1eW49dxD4NICQWgHcyF5p5Q+I
ZGlWkAjDRX0tGWZn3I/DRrHLmOt5BkV/si/WUH14ohjkqxNaCspdagYh5N2W8y7e
QZ7OIVkHRkA1huShj1si2QW3ecf17AaGT26DFwIDAQABAoIBAFx8lrIkKx/kOE0F
nY7BE9zkWGT5pFbsJpccoyqVW7iqEazoedK92+iI9J7aeN3l4vslSMAm0tN8hcg6
PWiENh/QcneyHh7jHZC8sqkfYPWdZmwaeOh2d1g8PFkll48lpPQEEM1CJa7PsPED
vsCdSTlXu+Yaan5PwOqQc+F+q8ZaYFk7lmCEgOWVTgwA6M7DOp7rcEzTnuJY9cie
PfHOZRrPOT+/V01hiKsbUWYvb2MWbfF/wxJwRAsG4u7Z6sQUDYipcQ85CYxt6m+t
YOap2Az+MSnGeM8Pf+KA2EKoR9Q7JejcD056DVAJcMuQNDguWxz6nje0vWRlL79l
Haym5gkCgYEAw9AY88DomONfnLdYNNi12F5wQKKg4Uz46T+S0SK6KOrHs7J57VvW
4Kpy+jqHTDxpeCp2IOECBPT5nyMyfLsyBlFAVjw0b+i6STg1uw/8xyl8+Dma0v36
DfUwiCrCmRSTzlHWf3cKNtjheGNlFy8qIfWFBznIpnMZX3lgxDBhW4UCgYEA2Z/l
N8zIn7SNAVWcSCYuk2sy5LgkK4pozx9BXzcBlXai/EgD65DVYylQaEfl5QOTBE3K
et3iGGPJRAakW6pC6J9RsRO6HfGJdkuOeVSAKgm36Tm4rYfV0bu7hH+xssD0d9fd
2W8RynFtnkaaov9mnbIp43JHcVl0WnlVnFApgOsCgYEApeYYTeSB7I6vghJgXB3D
K3cPueNPVLMnLE8db6zxdhs8aQXsgWpPCne/BDw0RyXj4dhvzvl0AYkgOHDUpJLh
FjMexDEr6CiQM9q4wy0PaBnBdHkxsFNX2R2EKcm4p4OkmqgBiGrtr3xewuXLTzI5
ix39wBp34nYf6CDpGC85PRUCgYEArA40hSdMvqdai+GJi6lUTY0FUbscLahiMM7/
Oi4c/HQta9Pr9YQukRWK0sd1RNjMlSyDlxxxsuLBrxypOSelepDrX1q/XQknqvUV
kWtzYMkKNEREdD3emNEZ8ima7j6LiWyLo2qi4DFJf0dG3vOZx7eiUoZ5YW5eBWHE
g68FAT0CgYAgA4wOsUTVcNUm2eG903h8ogjhmPRXwnEg6qwmFXdhD3vWUzFx214m
Tr50O1kW46VxSXjTP/yVxolcSrlFzB15sZk5dHUaNiTj+KU4lNAV5IMc02JP2EbM
RY0jOBEaIAvHslZDNtyIccVKFLDjKtutc9BB0V/YvZvysoY9UxVJAA==
-----END RSA PRIVATE KEY-----`
	privateKey, err := PemToPrivateKey(pri)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}
	signature, err := SignData(encrData, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("签名值:", signature)
}

func TestDESSign(t *testing.T) {
	datas := "gzDzraJo5CisyMhC5dSCaWjae3yTIXW6"
	privates := "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCJwlYeNsKdlRxbf8+SFteXsChSTKReMlGen++oxjYgN41H3963DcRb1EDAH1Ine/deyi4aFf108thuR9ocnZMV5BnPFWMp3+HJWTfY2MOdK/V9rfLUU6dssoM7MPt/izK67vvPeKy7Hv0+9VqOaQKURGbV/sjAFolgcbNH64b0DX9WF4K68fZG00J/w+5jkkXZidW8Z3E8a2NN2azlFYSlHYMG8ukpYTukV7ALvEKQIA9OCzjU24+os0QGfMzpjR51HC9dq4fIa+8U0abbFqeCoYRz16MomxSoXP9qj8P8w4lG1TxDT/0uqhUexGBPI+1yzpbvdXqNClu6OFVxdQHVAgMBAAECggEARQi4M0fp2yJAJbI1CNadj4xdiHCT5gh9Ump/pSo/MYHlMOGFMGKbxhDlqeGIP/Ulj8DtvXDLpPGVeB5VtQVaEhxOurHTEcb51Pb6v5ZQ4NCIo0SqbbDGM/h5Pw5a2h2dfIQKeHvWw6bR7dzyVm9VNYvZpN+bJnekvEn+N8pVxLLKNaa6dn0Aly+M6hlVUuKutIaynsK9qTef50AvrS487MHVwvt+gb6m/KET1Yf7t+ny/40b1h6DgIdkK/wejmxi9mw1ayuFb3NktBnOHJOjCWYjWGMa4VZ8NuIyBMT92drcGbJtO8xwL80Ln2Ps3hYC7ohGJL9E3av5sCHI9Z/IwQKBgQDJJfR2nAEUeTrEYInzqlHQr7hhl6MsJOdAhaY09ZrGyRahdfqSoMPr6pSpZuVA01QWLTqYtn+X5WTG7KqcehcsobXxKpyA8BceK3P328rQvxl9U/BvAGwz2BYuVAbrL+PgAZv0uFEY8gR5Epl7dz+9UWnBi9KwHa7mHw3N1TaNjQKBgQCvUzjjRCB0JJ3my2nwJu0KNTwMOe1QhP5F7iFcjJQW/oFtrlTaEheDnTNovuhNPVpRAwGlARQyA4uQY+rFJGdzbtoHk+EsJNLvsJV+Tq5yh2hv5ked4iyfA5RWd37FYmfKqM6nUyvSSoschQXtsr/nEosB99JLvrflcks2iEx/aQKBgD5h1Ag456jW1B/1JLN5/fevl4pEwek95K5BBMPl68N8t9UJRtXUoA55aPOEotLQ94INMuALsVSFYxTCb0MqJifEWy3ZHkJqs3C63zNeae8FZT1WG/oA8o29lVt22dJ0vsJJHXnu88+9tx9pYkpFOHJZXmgVGhlei1B5Dwnn9ww9AoGAVHQhNhBuFaRBz5fyqvUFP+KOz1DkCOJXXbYsqdkpyL3F+OB+DSGj5AlIZ092tSY1qEprc2FGqiTdCKuovlgf4RHnwriwQcRnO4BzMomSLKcfXq+tldcKKXre7JvZHBmf55ZTHXTJ6h1wT0egqHRvTk63WTZYPZZcHRFmO5mCR+kCgYEAkrKArJRg9VO9rVXg6pa1QZJ8NdOOrDx1MY0A+plPNww6ht+bh6c5WCb2nXmV97huGQIGZFgHdAEnLM+EjcX2/ZH5rPjpdP16G8n8TSwXrJOyeDj21odr1DmHw7cQb1deibaQ9eK9DROZEnO/XV4jJtYKomKUcTH99JMm8dUNzIM="
	sign := DesSignData(datas, privates)
	fmt.Println(sign)
}
