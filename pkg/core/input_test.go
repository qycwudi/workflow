package core

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestParseNodeInputs(t *testing.T) {
	executionContext := NewExecutionContext(context.Background(), "test", "test", nil)
	executionContext.SetVariable("start_0.output", `{"start_0":{"Input":{"sssss":"我是字符串","oooooo":{"ssssssssssss":"我是对象"},"arr_str":["我是数组字符串-1","我是数组字符串-2"],"array_obj":[{"int":18,"str":"我是数组里的对象的字符串"}],"enable":true,"query":"Hello Flow."}`)
	inputValues := map[string]NodeDataInputsValues{}
	ivjs := `{
                    "sss": {
                        "type": "ref",
                        "content": [
                            "start_0",
                            "oooooo",
                            "ssssssssssss"
                        ]
                    },
                    "result": {
                        "type": "ref",
                        "content": [
                            "start_0",
                            "query"
                        ]
                    },
                    "arrrrrr": {
                        "type": "ref",
                        "content": [
                            "start_0",
                            "arr_str"
                        ]
                    },
                    "sssssss": {
                        "type": "ref",
                        "content": [
                            "start_0",
                            "oooooo",
                            "ssssssssssss"
                        ]
                    }
                }`
	err := json.Unmarshal([]byte(ivjs), &inputValues)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	ijs := `{
                    "type": "object",
                    "properties": {
                        "sss": {
                            "type": "string",
                            "items": {
                                "type": "string"
                            }
                        },
                        "result": {
                            "type": "string"
                        },
                        "arrrrrr": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        },
                        "sssssss": {
                            "type": "string"
                        }
                    }
                }`

	var inputs NodeDataInputs
	err = json.Unmarshal([]byte(ijs), &inputs)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	got, err := ParseNodeInputs(executionContext, inputValues, inputs)
	if err != nil {
		t.Fatalf("ParseNodeInputs failed: %v", err)
	}
	logx.Debugf("[输入处理] got:%+v", got)
}

func TestProcessNodeOutput(t *testing.T) {
	jsonData := `{\\\"start_0\\\":{\\\"Input\\\":{\\\"sssss\\\":\\\"我是字符串\\\",\\\"oooooo\\\":{\\\"ssssssssssss\\\":\\\"我是对象\\\"},\\\"arr_str\\\":[\\\"我是数组字符串-1\\\",\\\"我是数组字符串-2\\\"],\\\"array_obj\\\":[{\\\"int\\\":18,\\\"str\\\":\\\"我是数组里的对象的字符串\\\"}],\\\"enable\\\":true,\\\"query\\\":\\\"Hello Flow.\\\"}}`
	gjson.Get(jsonData, "start_0.Input.sssss")
	logx.Debugf("[输出处理] jsonData:%s", jsonData)
}

func TestConvertValue(t *testing.T) {
	j := "1"
	id, err := json.Marshal(j)
	t.Logf("id:%s,err:%v", id, err)
}
