package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestDESEncrypt(t *testing.T) {
	// 初始化随机数种子
	name := "张三"
	idCard := "430721199409256116"
	startDate := "20130819"
	endDate := "20230819"

	// 1. 封装请求数据
	condition, err := dataEncapsulation(name, idCard, startDate, endDate)
	if err != nil {
		fmt.Printf("Error encapsulating data: %v\n", err)
		return
	}
	fmt.Println("---- Request Parameters ----")
	// 打印 map 内容以便调试 (类似 Java 的 JSONUtil.toJsonStr(multiValueMap))
	mapBytes, _ := json.MarshalIndent(condition, "", "  ")
	fmt.Println(string(mapBytes))
	fmt.Println("----")

	// 2. 发送 HTTP POST 请求 (Form-UrlEncoded)
	formData := url.Values{}
	for key, value := range condition {
		// 确保所有值都作为字符串发送
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	// 创建带超时的 HTTP Client
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置 10 秒超时，对应 Java 代码中的 10000ms
	}

	resp, err := client.PostForm(apiURL, formData)
	if err != nil {
		fmt.Printf("Error making POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 3. 读取并打印响应
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println("---- API Response ----")
	fmt.Println(string(bodyBytes))
}

// --- 配置常量 ---
const (
	userCode   = "HFJK20250414180959"
	appCode    = "HFJKAPP20250414181054"
	appKey     = "1tWKj3UmkoxhL0nqFh9/JT5M"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 // 3DES 密钥
	privateKey = `MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCEV6drAbzSkymudXog5FVXFkjHk+gjLqDqDQozvHMUkPxk7g5MtJ2rqFGmLlao1nnk7PNa6C35+b0Eo4smRsczgaZBg3ys7PNkkrjgdFlvRitt0epy5+s1r7MQdyIQCcnJwXzgUPwNVn3RB61bbYwZheWIQ5DUdT+xFNUhzhH1f65RZWLJ6my7YNur07WJ0ask1dn8aMcol8tg1C7ZXnSYsLevItfwu7HXw3lJ8bihkSr1H7yGjAzitdNbYMhUG/ZHn//nWFH0kabT585JI10I/tVtaOvqS47UmcL9Jx4fZvyx2CQNOTFB2i3UjpXOnJ8yL0hAYbx7R2civsSQiTb5AgMBAAECggEAL9YjPOmm4BKzqUzrdUIzmsQCUKmk6jqrjY5jFqpSiqeRo8Xsw0syBt9TjBcJ2gOUkewYVs6/97CO40NeQ9qnnHWLq/ikMzl/DOaADxnfTfp2Lj8HWLt96Kz2s8fsNxHir5vR0J9VTFjsJ4d34Slqz7r3adbKXhF3kUGvfIWpNiyB7j3yOyiSC2NE2v6gGa9i83El3NEPSTyZFP22aZZhYHx6plgk80lCI1grLjXR51pWizXvhsu/giGoWxBnpygIWaQJOmsP5q9816f2npiiH02/M4sCtExJMIeOv23borpzGU/POFHnfO6M4BGNeeCGSNyVO1O0Tj+Fn1tChFyBgQKBgQD9Fqk7bYgFrb/jhBcRBnQANopUJEpXAodGvm7elS/6pG1sENaxh7BCou4tiW2xZikxaNUTwkaJD1HkVaw0wAQHbuVNxIa/MCLN0pKxSBQwkG5+DWR9OOrf2nDcc9CgCGLLqTFVYkBHiyNUZr3Y0Nk2oUJYA1d3HHzzfXzG1pxpMQKBgQCF3WYYSM4GOcKN0H06VXx1vfon5uR7CeSICDf/JeG9BaF8Mkaf58tTDyhfRo/kjS2Ig8gVrfaqe407TIvAdqYL5gc1bruLuVzcD9vZKhbvm7guH2xhIIvcoG83fVjAHM7O7ByxHvWWKHcXhIE8/fU0OZ6fCa/2nI2Rc6U9wpS4SQKBgAH/HMSoP4tz7HCaCSiMOXwK8hpp0uBO18xaEjvVR3SduXNByr/Jxz0vVdusGt5sZOTL4+ps/Ot14HqBpoMgBKgeWta7Nyjo801WXOvjGq2rZgO+jU1jlo6/hWZnz3yDtGvQ4N/Aj6tA0L2ItfSw6RXCPei91gHNirrNAZ/7723xAoGAATN3Uwh0MMIb6MHSHc/eif/mWq1Lp37zOfU462UfnV6LkF2zHIQr4tbj+dbcO6S4n9zu0qL475akMcACEPv/iWPK9MPFkv9awe6yfkROaF/xSxilFEoCdzxJQyowvaiEgn7D0yL/+RTr3J4nurBcntdVgP/JQGEvM/rhpKg2RWECgYAx4JiP8/rPPLM7BxiTkXF8YSm0xM2epSp8SdCnag36AC2X2eoj4HLg1KzAHoLy8jLpo26lsrxdooUIhvQVnD76fRAT3H2keRZR3a45g8sK+6G699EGpG/3PrbNtZlZDFU5oHawZNIxrLTD7ZC+EfNtLKNTflWq+nHqk8DK6NgDWw==` // RSA 私钥
	// publicKey = "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAJJY/kusrKic6XAhHMakIX06nEnMnvceRV5m8dLIKQcz0RNkDXiEQ/EIV0hZQNQlTIyB5f6OzQeDjJlDUGgkyD8CAwEAAQ==" // 公钥在此示例中未使用
	apiURL = "https://api.njzhdkj.com/superapi/super/idcard/validity"
)

// --- JSON 结构体定义 ---

// RequestHeader 对应 Java 代码中的 headerJson
type RequestHeader struct {
	QryBatchNo string `json:"qryBatchNo"` // 验证批次号
	UserCode   string `json:"userCode"`   // 商户编号
	SysCode    string `json:"sysCode"`    // 应用编号
	QryReason  string `json:"qryReason"`  // 原因
	QryDate    string `json:"qryDate"`    // 格式：yyyyMMdd
	QryTime    string `json:"qryTime"`    // 格式：hhmmss
	Version    string `json:"version"`    // 版本号
}

// RequestCondition 对应 Java 代码中的 conditionJson
type RequestCondition struct {
	RealName  string `json:"realName"`
	IdCard    string `json:"idCard"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	// Photo string `json:"photo,omitempty"` // 如果需要照片字段，取消注释
}

// RequestPayload 对应 Java 代码中的 allJson
type RequestPayload struct {
	Header    RequestHeader    `json:"header"`
	Condition RequestCondition `json:"condition"`
}

// dataEncapsulation 封装请求数据，返回用于 POST 请求的 map 和可能的错误
func dataEncapsulation(name, idCard, startDate, endDate string) (map[string]interface{}, error) {

	// 获取日期 yyyyMMdd
	qryDate := "20250415"

	// 获取时间 HHmmss
	qryTime := "183209"

	// // 生成商户批次号：14位时间戳 + 4位随机数 (Java代码是*10000，产生0-9999，所以是4位)
	// timestamp := now.Format("20060102150405")
	// random := rand.Intn(10000)                             // 生成 0 到 9999 之间的随机数
	qryBatchNo := "202504151810269831" // 使用 %04d 保证随机数部分总是4位
	fmt.Println("商户批次号：" + qryBatchNo)

	// --- 生成请求报文 ---
	header := RequestHeader{
		QryBatchNo: qryBatchNo,
		UserCode:   userCode,
		SysCode:    appCode,
		QryReason:  "查询", // 与 Java 代码一致
		QryDate:    qryDate,
		QryTime:    qryTime,
		Version:    "2.0", // 与 Java 代码一致
	}

	condition := RequestCondition{
		RealName:  name,
		IdCard:    idCard,
		StartDate: startDate,
		EndDate:   endDate,
		// Photo: photo, // 如果需要
	}

	payload := RequestPayload{
		Header:    header,
		Condition: condition,
	}

	// 将整个 payload 结构体序列化为 JSON 字符串
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}
	dataCompress := string(payloadBytes)

	// --- 加密请求报文 ---
	// 注意：Java 代码硬编码了 IV "1234567a"。Go 工具函数也需要这个 IV。
	iv := "1234567a"
	encrData := Encrypt3DES(dataCompress, appKey, iv)
	// 注意：提供的 Go utils 函数在出错时直接返回错误字符串，而不是 error 类型。
	// 在实际应用中，最好修改 utils 函数返回 (string, error)。
	// 这里我们假设如果加密失败，encrData 会包含错误信息。
	// 更好的检查方式可能是检查返回的字符串是否像 Base64 编码。
	fmt.Println("加密报文：" + encrData)
	// (可选) 添加更健壮的错误检查，例如检查 encrData 是否为空或包含特定错误关键字
	if encrData == "" || strings.Contains(encrData, "error") || strings.Contains(encrData, "failed") {
		// 尝试从返回的字符串判断是否是错误信息
		if strings.HasPrefix(encrData, "failed") || strings.HasPrefix(encrData, "key length must be") {
			return nil, fmt.Errorf("DES encryption failed: %s", encrData)
		}
		// 如果不能确定，可以加一个通用错误
		// return nil, fmt.Errorf("DES encryption might have failed, returned: %s", encrData)
	}

	// --- 对加密后的报文进行签名 ---
	// 注意：utils.DESSign 函数名有误导性，它实际执行的是 RSA 签名
	signature := SignRSA_MD5(encrData, privateKey)
	fmt.Println("签名值：" + signature)
	// (可选) 添加更健壮的错误检查
	if signature == "" || strings.Contains(signature, "error") || strings.Contains(signature, "failed") {
		// 尝试从返回的字符串判断是否是错误信息
		if strings.HasPrefix(signature, "failed") {
			return nil, fmt.Errorf("RSA signing failed: %s", signature)
		}
		// 如果不能确定，可以加一个通用错误
		// return nil, fmt.Errorf("RSA signing might have failed, returned: %s", signature)
	}

	// --- 组装最终请求参数 ---
	requestParams := map[string]interface{}{
		"condition": encrData,
		"userCode":  userCode,
		"signature": signature,
		"vector":    iv, // 将 IV 也添加到请求参数中
	}

	return requestParams, nil
}
