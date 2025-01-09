package rulego

import (
	"fmt"
	"testing"

	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_ValueSubstitution(t *testing.T) {
	msg := map[string]any{
		"intA": 20,
		"intB": 30,
		"user": map[string]any{
			"name": "张三",
			"age":  25,
			"address": map[string]any{
				"city":   "北京",
				"street": "朝阳路",
			},
			"hobbies": []string{"游泳", "跑步", "阅读"},
			"scores":  []int{90, 85, 95},
		},
	}

	param := `
	<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
	<Add xmlns="http://tempuri.org/">
	<intA>${intA}</intA>
	<intB>${intB}</intB>
	<intC>${intC}</intC>
	<intD>12</intD>
	<user>
		<name>${user.name}</name>
		<age>${user.age}</age>
		<address>
			<city>${user.address.city}</city>
			<street>${user.address.street}</street>
		</address>
		<hobbies>
			<hobby>${user.hobbies[0]}</hobby>
			<hobby>${user.hobbies[1]}</hobby>
			<hobby>${user.hobbies[2]}</hobby>
		</hobbies>
		<scores>
			<score>${user.scores}</score>
		</scores>
	</user>
	</Add>
	</soap:Body>
	</soap:Envelope>
	`

	jsonData, _ := json.Marshal(msg)
	logx.Info(string(jsonData))
	result := replaceXmlTemplateVars(param, string(jsonData))

	// 期望的结果
	expected := `
	<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
	<Add xmlns="http://tempuri.org/">
	<intA>20</intA>
	<intB>30</intB>
	<intC></intC>
	<intD>12</intD>
	<user>
		<name>张三</name>
		<age>25</age>
		<address>
			<city>北京</city>
			<street>朝阳路</street>
		</address>
		<hobbies>
			<hobby>游泳</hobby>
			<hobby>跑步</hobby>
			<hobby>阅读</hobby>
		</hobbies>
		<scores>
			<score>90,85,95</score>
		</scores>
	</user>
	</Add>
	</soap:Body>
	</soap:Envelope>
	`

	if result != expected {
		t.Errorf("Value substitution failed.\nExpected:\n%s\nGot:\n%s", expected, result)
	}

	fmt.Printf("Mapped XML:\n%s", result)
}



func Test_CZ(t *testing.T) {
	msg := map[string]any{
		"intA": 20,
		"intB": 30,
		"user": map[string]any{
			"name": "张三",
			"age":  25,
			"address": map[string]any{
				"city":   "北京",
				"street": "朝阳路",
			},
			"hobbies": []string{"游泳", "跑步", "阅读"},
			"scores":  []int{90, 85, 95},
		},
	}

	param := `
	<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
	<Add xmlns="http://tempuri.org/">
	<intA>${intA}</intA>
	<intB>${intB}</intB>
	<intC>${intC}</intC>
	<intD>12</intD>
	<user>
		<name>${user.name}</name>
		<age>${user.age}</age>
		<address>
			<city>${user.address.city}</city>
			<street>${user.address.street}</street>
		</address>
		<hobbies>
			<hobby>${user.hobbies[0]}</hobby>
			<hobby>${user.hobbies[1]}</hobby>
			<hobby>${user.hobbies[2]}</hobby>
		</hobbies>
		<scores>
			<score>${user.scores}</score>
		</scores>
	</user>
	</Add>
	</soap:Body>
	</soap:Envelope>
	`

	jsonData, _ := json.Marshal(msg)
	logx.Info(string(jsonData))
	result := replaceXmlTemplateVars(param, string(jsonData))


	fmt.Printf("Mapped XML:\n%s", result)
}