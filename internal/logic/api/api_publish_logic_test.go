package api

import (
	"os"
	"testing"
)

func Test_getApiDoc(t *testing.T) {

	markdown, err := GenerateMarkdownFromDSL(dsl, "http://127.0.0.1:8888", "123", "namenamename", "这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口这是一个测试接口")
	if err != nil {
		t.Errorf("GenerateMarkdownFromDSL() error = %v", err)
	}
	t.Log(markdown)
	// 写入文件
	err = os.WriteFile("api_doc.md", []byte(markdown), 0644)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}
}

const dsl = `
{
  "edges": [
    {
      "sourceNodeID": "start_0",
      "targetNodeID": "llm_1750924136981"
    },
    {
      "sourceNodeID": "llm_1750924136981",
      "targetNodeID": "end_0"
    }
  ],
  "nodes": [
    {
      "id": "start_0",
      "data": {
        "title": "开始",
        "outputs": {
          "type": "object",
          "required": [
            "input",
            "obj",
            "bl",
            "arr_str"
          ],
          "properties": {
            "bl": {
              "key": 16,
              "name": "bl",
              "type": "boolean",
              "extra": {
                "index": 2
              },
              "default": true,
              "description": "布尔值类型的",
              "isPropertyRequired": true
            },
            "obj": {
              "key": 15,
              "name": "obj",
              "type": "object",
              "extra": {
                "index": 1
              },
              "default": "{\n    \"name\": \"ssssss\",\n    \"age\": 18\n}",
              "required": [],
              "properties": {
                "name": {
                  "key": 18,
                  "name": "name",
                  "type": "string",
                  "extra": {
                    "index": 1
                  },
                  "description": "姓名",
                  "isPropertyRequired": false
                }
              },
              "description": "这是一个标准的 json",
              "isPropertyRequired": true
            },
            "input": {
              "key": 4,
              "name": "input",
              "type": "string",
              "extra": {
                "index": 0
              },
              "default": "帮我写一个生日祝福语",
              "description": "大模型输入",
              "isPropertyRequired": true
            },
            "arr_str": {
              "key": 17,
              "name": "arr_str",
              "type": "array",
              "extra": {
                "index": 3
              },
              "items": {
                "type": "string"
              },
              "description": "字符串数组",
              "isPropertyRequired": true
            }
          }
        }
      },
      "meta": {
        "position": {
          "x": -15,
          "y": -142
        }
      },
      "type": "start"
    },
    {
      "id": "end_0",
      "data": {
        "title": "结束",
        "outputs": {
          "type": "object",
          "properties": {
            "data": {
              "type": "string",
              "description": "这是一个大模型数据数据的参数"
            },
            "content": {
              "type": "string",
              "description": "这是一个大模型生成文本结果"
            }
          }
        },
        "inputsValues": {
          "data": {
            "type": "ref",
            "content": [
              "llm_1750924136981",
              "content"
            ]
          },
          "content": {
            "type": "ref",
            "content": [
              "llm_1750924136981",
              "content"
            ]
          }
        }
      },
      "meta": {
        "position": {
          "x": 755,
          "y": -142
        }
      },
      "type": "end"
    },
    {
      "id": "llm_1750924136981",
      "data": {
        "title": "大模型调用",
        "custom": {
          "retry": 3,
          "tools": [],
          "output": {
            "type": "object",
            "required": [
              "content",
              "totalTokens",
              "promptTokens",
              "completionTokens"
            ],
            "properties": {
              "content": {
                "type": "string",
                "default": ""
              },
              "totalTokens": {
                "type": "integer",
                "default": 0
              },
              "promptTokens": {
                "type": "integer",
                "default": 0
              },
              "completionTokens": {
                "type": "integer",
                "default": 0
              }
            }
          },
          "modelId": 3,
          "timeout": 30,
          "outputType": "string",
          "userPrompt": "{{input}}",
          "systemPrompt": "你是一个专业的AI助手，请根据用户的需求提供帮助。",
          "errorHandlingMode": "abort"
        },
        "inputs": {
          "properties": {
            "input": {
              "type": "string",
              "title": "input"
            }
          }
        },
        "outputs": {
          "type": "object",
          "required": [
            "content",
            "totalTokens",
            "promptTokens",
            "completionTokens"
          ],
          "properties": {
            "content": {
              "type": "string",
              "default": ""
            },
            "totalTokens": {
              "type": "integer",
              "default": 0
            },
            "promptTokens": {
              "type": "integer",
              "default": 0
            },
            "completionTokens": {
              "type": "integer",
              "default": 0
            }
          }
        },
        "inputsValues": {
          "input": {
            "type": "ref",
            "content": [
              "start_0",
              "input"
            ]
          }
        }
      },
      "meta": {
        "position": {
          "x": 368,
          "y": 50
        }
      },
      "type": "llm"
    }
  ]
}`
