package workflow

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

func TestRegister(t *testing.T) {
	dsl := `
	{
    "edges": [
        {
            "sourceNodeID": "start_0",
            "targetNodeID": "end_0"
        }
    ],
    "nodes": [
        {
            "id": "start_0",
            "data": {
                "title": "Start",
                "outputs": {
                    "type": "object",
                    "required": [
                        "array_obj",
                        "enable",
                        "query",
                        "sssss"
                    ],
                    "properties": {
                        "query": {
                            "key": 12,
                            "name": "query",
                            "type": "string",
                            "extra": {
                                "index": 2
                            },
                            "default": "Hello Flow.",
                            "isPropertyRequired": true
                        },
                        "sssss": {
                            "key": 15,
                            "name": "sssss",
                            "type": "string",
                            "extra": {
                                "index": 3
                            },
                            "isPropertyRequired": true
                        },
                        "enable": {
                            "key": 11,
                            "name": "enable",
                            "type": "boolean",
                            "extra": {
                                "index": 1
                            },
                            "default": true,
                            "isPropertyRequired": true
                        },
                        "oooooo": {
                            "key": 41,
                            "name": "oooooo",
                            "type": "object",
                            "extra": {
                                "index": 4
                            },
                            "required": [],
                            "properties": {
                                "ssssssssssss": {
                                    "key": 42,
                                    "name": "ssssssssssss",
                                    "type": "string",
                                    "extra": {
                                        "index": 1
                                    }
                                }
                            },
                            "isPropertyRequired": false
                        },
                        "arr_str": {
                            "key": 51,
                            "name": "arr_str",
                            "type": "array",
                            "extra": {
                                "index": 6
                            },
                            "items": {
                                "type": "string"
                            }
                        },
                        "array_obj": {
                            "key": 10,
                            "name": "array_obj",
                            "type": "array",
                            "extra": {
                                "index": 0
                            },
                            "items": {
                                "type": "object",
                                "required": [
                                    "int",
                                    "str"
                                ],
                                "properties": {
                                    "int": {
                                        "key": 39,
                                        "name": "int",
                                        "type": "number",
                                        "extra": {
                                            "index": 0
                                        },
                                        "isPropertyRequired": true
                                    },
                                    "str": {
                                        "key": 40,
                                        "name": "str",
                                        "type": "string",
                                        "extra": {
                                            "index": 1
                                        },
                                        "isPropertyRequired": true
                                    }
                                }
                            },
                            "isPropertyRequired": true
                        }
                    }
                }
            },
            "meta": {
                "position": {
                    "x": 259.21613278183247,
                    "y": 78.02894601785266
                }
            },
            "type": "start"
        },
        {
            "id": "end_0",
            "data": {
                "title": "End",
                "outputs": {
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
                },
                "inputsValues": {
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
                }
            },
            "meta": {
                "position": {
                    "x": 911.2161327818328,
                    "y": 104.02894601785266
                }
            },
            "type": "end"
        }
    ]
}
	`

	var workflowDefinition core.WorkflowDef
	err := sonic.Unmarshal([]byte(dsl), &workflowDefinition)
	if err != nil {
		logx.Errorf("[工作流] 解析工作流文件失败 [错误:%s]", err)
	}
	definitionBytes, _ := sonic.MarshalIndent(workflowDefinition, "", "  ")
	logx.Debugf("[工作流] 工作流DSL [ID:%s] [DSL:%s]", "test", string(definitionBytes))
}
