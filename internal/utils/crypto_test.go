package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
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

func TestAesEncrypt2(t *testing.T) {
	key := "vQRGdmlvZ@KP$o7Z"
	req := []map[string]interface{}{
		{
			"authorize": 1,
			"id":        "420116199912165234",
			"name":      "李欣瑞",
		},
	}
	js, _ := json.Marshal(req)
	encrypt := FyyAesEncrypt(string(js), key)
	fmt.Println(encrypt)
	// hAzC8jMn2ehi50dlKhroeqtxTbGFJqkCukNIcLX62FVc/bcj44RAhUUJKMhT1RoExgZZO9rs1M45DmoWkr3nXQ==
	// OwiCTaEemLZZfbLEtacqur7SAtXB18tjxaT2+b5dLi5gQqd0imG4Qt2s7FF2xkkKu5k64zxzoYjTbvEKqQvYis2QLzxvn804EG0vOIrd70MzRicGL4khUtfYZ3zLR75br86cze+FkLrZ
}

func TestAesDecrypt2(t *testing.T) {
	key := "vQRGdmlvZ@KP$o7Z"
	content := "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"
	decrypt := FyyAesDecrypt(content, key)
	fmt.Println(decrypt)
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
	// hAzC8jMn2ehi50dlKhroeqtxTbGFJqkCukNIcLX62FVc/bcj44RAhUUJKMhT1RoExgZZO9rs1M45DmoWkr3nXQ==
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
