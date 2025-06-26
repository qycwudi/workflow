package api

import (
	"context"
	errors2 "errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/dispatch/broadcast"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiPublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPublishLogic {
	return &ApiPublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiPublishLogic) ApiPublish(req *types.ApiPublishRequest) (resp *types.ApiPublishResponse, err error) {
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 检查画布名称重复
	_, err = l.svcCtx.ApiModel.FindByName(l.ctx, req.ApiName)
	if !errors2.Is(err, sqlc.ErrNotFound) {
		return nil, errors.New(int(logic.SystemStoreError), "API 名称重复")
	}
	// 自动保存一个历史版本
	history, err := l.svcCtx.CanvasHistoryModel.Insert(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.Id,
		Draft:       canvas.Draft,
		Name:        req.ApiName,
		CreateTime:  time.Now(),
		Mode:        model.CanvasHistoryModeApi,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "保存历史版本失败")
	}
	historyId, err := history.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取历史版本ID失败")
	}

	// 查询有没有发布过api
	api, err := l.svcCtx.ApiModel.FindByWorkspaceId(l.ctx, req.Id)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemStoreError), "查询API失败")
	}
	var apiId string
	tagJson, err := sonic.Marshal(req.Tag)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "标签转换失败")
	}
	if api == nil {
		apiId = xid.New().String()
		_, err = l.svcCtx.ApiModel.Insert(l.ctx, &model.Api{
			WorkspaceId: req.Id,
			ApiId:       apiId,
			ApiName:     req.ApiName,
			ApiDesc:     req.ApiDesc,
			ApiDoc:      "",
			Tag:         string(tagJson),
			Dsl:         canvas.Draft,
			Status:      model.ApiStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "发布 API 失败")
		}
	} else {
		apiId = api.ApiId
		// 如果发布过，则更新
		err = l.svcCtx.ApiModel.Update(l.ctx, &model.Api{
			Id:          api.Id,
			WorkspaceId: api.WorkspaceId,
			ApiId:       api.ApiId,
			ApiName:     req.ApiName,
			ApiDesc:     req.ApiDesc,
			Tag:         string(tagJson),
			Dsl:         canvas.Draft,
			Status:      model.ApiStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  api.CreateTime,
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "更新 API 失败")
		}
	}

	// 3. 发送加载链服务消息
	err = broadcast.NewApiLoadSync().Publish(l.ctx, &broadcast.ApiLoadSyncMsg{
		ApiId:     apiId,
		RuleChain: canvas.Draft,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
	}
	// 删除redis缓存
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.EnvRedisKey, apiId))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除API环境变量缓存失败")
	}
	// 更新 API 文档
	getApiDoc(l.ctx, l.svcCtx.ApiModel, l.svcCtx.Config.ApiUrl, apiId)
	resp = &types.ApiPublishResponse{ApiId: apiId}
	return resp, nil
}

func getApiDoc(ctx context.Context, apiModel model.ApiModel, apiUrl, apiId string) error {
	// 查询
	apiEntity, err := apiModel.FindOneByApiId(ctx, apiId)
	if err != nil {
		return err
	}
	apiDoc, err := GenerateMarkdownFromDSL(apiEntity.Dsl, apiUrl, apiId, apiEntity.ApiName, apiEntity.ApiDesc)
	if err != nil {
		return err
	}
	apiEntity.ApiDoc = apiDoc
	err = apiModel.Update(ctx, apiEntity)
	if err != nil {
		return err
	}
	return nil
}

// 生成 Markdown 文档的主函数
func GenerateMarkdownFromDSL(dslJSON, apiUrl, apiId, name, desc string) (string, error) {
	// 查找 start_0 和 end_0 节点
	var startNodeData, endNodeData gjson.Result
	nodes := gjson.Get(dslJSON, "nodes")

	nodes.ForEach(func(key, value gjson.Result) bool {
		nodeId := value.Get("id").String()
		if strings.HasPrefix(nodeId, "start_") {
			startNodeData = value
		} else if strings.HasPrefix(nodeId, "end_") {
			endNodeData = value
		}
		return true
	})

	if !startNodeData.Exists() {
		return "", fmt.Errorf("未找到 start_0 节点")
	}

	// 生成简化的 API 文档
	markdown := generateSimpleApiDoc(startNodeData, endNodeData, apiUrl, apiId, name, desc)
	return markdown, nil
}

// 生成简化的API文档
func generateSimpleApiDoc(startNodeData, endNodeData gjson.Result, apiUrl, apiId, name, desc string) string {
	var md strings.Builder

	// 获取基本信息

	// 1. 接口请求地址
	md.WriteString("# API 接口文档\n\n")
	md.WriteString("## 接口信息\n\n")
	md.WriteString(fmt.Sprintf("- **接口名称**: %s\n", name))
	md.WriteString(fmt.Sprintf("- **接口描述**: %s\n", desc))
	md.WriteString(fmt.Sprintf("- **请求地址**: `%s/%s`\n", apiUrl, apiId))
	md.WriteString("- **请求方式**: `POST`\n")
	md.WriteString("- **Content-Type**: `application/json`\n\n")
	md.WriteString("- **Authorization**: `token(在密钥管理创建)`\n\n")

	// 2. 接口请求参数表格
	outputs := startNodeData.Get("data.outputs")
	properties := outputs.Get("properties")
	required := outputs.Get("required")

	md.WriteString("## 请求参数\n\n")

	// 收集并排序参数
	type ParamInfo struct {
		Name  string
		Data  gjson.Result
		Index int
	}

	var params []ParamInfo
	properties.ForEach(func(key, value gjson.Result) bool {
		paramName := key.String()
		index := value.Get("extra.index").Int()
		params = append(params, ParamInfo{
			Name:  paramName,
			Data:  value,
			Index: int(index),
		})
		return true
	})

	// 按 index 排序
	sort.Slice(params, func(i, j int) bool {
		return params[i].Index < params[j].Index
	})

	if len(params) == 0 {
		md.WriteString("无请求参数\n\n")
	} else {
		// 生成参数表格
		md.WriteString("| 参数名 | 类型 | 必填 | 默认值 | 说明 |\n")
		md.WriteString("|--------|------|------|--------|------|\n")

		for _, param := range params {
			name := param.Name
			prop := param.Data
			paramType := prop.Get("type").String()
			description := prop.Get("description").String()
			if description == "" {
				description = "-"
			}
			isReq := isRequiredParam(required, name)
			requiredText := map[bool]string{true: "是", false: "否"}[isReq]

			// 默认值
			defaultValue := "-"
			if prop.Get("default").Exists() {
				rawValue := fmt.Sprintf("%v", prop.Get("default").Value())
				// 删除换行符和其他可能破坏表格的字符
				cleanValue := strings.ReplaceAll(rawValue, "\n", "")
				cleanValue = strings.ReplaceAll(cleanValue, "\r", "")
				cleanValue = strings.ReplaceAll(cleanValue, "|", "\\|")
				defaultValue = fmt.Sprintf("`%s`", cleanValue)
			}

			// 转义特殊字符防止表格破坏
			safeName := strings.ReplaceAll(name, "|", "\\|")
			safeParamType := strings.ReplaceAll(paramType, "|", "\\|")
			safeDescription := strings.ReplaceAll(description, "|", "\\|")

			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				safeName, safeParamType, requiredText, defaultValue, safeDescription))
		}
		md.WriteString("\n")

		// 单独处理对象属性表格
		for _, param := range params {
			name := param.Name
			prop := param.Data
			paramType := prop.Get("type").String()

			if paramType == "object" {
				objectProps := prop.Get("properties")
				if objectProps.Exists() {
					md.WriteString("**" + name + " 对象属性：**\n\n")
					md.WriteString("| 属性名 | 类型 | 说明 |\n")
					md.WriteString("|--------|------|------|\n")
					objectProps.ForEach(func(key, value gjson.Result) bool {
						propName := key.String()
						propType := value.Get("type").String()
						propDesc := value.Get("description").String()
						if propDesc == "" {
							propDesc = "-"
						}
						// 转义特殊字符防止表格破坏
						safePropName := strings.ReplaceAll(propName, "|", "\\|")
						safePropType := strings.ReplaceAll(propType, "|", "\\|")
						safePropDesc := strings.ReplaceAll(propDesc, "|", "\\|")

						md.WriteString(fmt.Sprintf("| %s | %s | %s |\n", safePropName, safePropType, safePropDesc))
						return true
					})
					md.WriteString("\n")
				}
			}
		}
	}

	// 3. 请求示例
	md.WriteString("## 请求示例\n\n")
	md.WriteString("```bash\n")
	md.WriteString(fmt.Sprintf("curl -X POST %s/%s \\\n", apiUrl, apiId))
	md.WriteString("  -H 'Content-Type: application/json' \\\n")
	md.WriteString("  -d '")

	// 生成请求参数示例
	md.WriteString("{\n")
	for i, param := range params {
		name := param.Name
		prop := param.Data

		if i > 0 {
			md.WriteString(",\n")
		}

		example := generateExampleValue(prop)
		md.WriteString(fmt.Sprintf("    \"%s\": %s", name, example))
	}
	md.WriteString("\n  }'\n```\n\n")

	// JSON 格式示例
	md.WriteString("**JSON 格式：**\n\n")
	md.WriteString("```json\n")
	md.WriteString("{\n")
	for i, param := range params {
		name := param.Name
		prop := param.Data

		if i > 0 {
			md.WriteString(",\n")
		}

		example := generateExampleValue(prop)
		md.WriteString(fmt.Sprintf("  \"%s\": %s", name, example))
	}
	md.WriteString("\n}\n```\n\n")

	// 4. 返回值说明
	if endNodeData.Exists() {
		endOutputs := endNodeData.Get("data.outputs")
		endProperties := endOutputs.Get("properties")

		md.WriteString("## 返回值\n\n")

		if !endProperties.Exists() || endProperties.Map() == nil || len(endProperties.Map()) == 0 {
			md.WriteString("无返回值\n\n")
		} else {
			// 统计返回值字段数量
			returnFieldCount := 0
			endProperties.ForEach(func(key, value gjson.Result) bool {
				returnFieldCount++
				return true
			})

			md.WriteString(fmt.Sprintf("接口返回 **%d** 个字段：\n\n", returnFieldCount))

			// 生成返回值表格
			md.WriteString("| 字段名 | 类型 | 说明 |\n")
			md.WriteString("|--------|------|------|\n")

			endProperties.ForEach(func(key, value gjson.Result) bool {
				fieldName := key.String()
				fieldType := value.Get("type").String()
				fieldDesc := value.Get("description").String()
				if fieldDesc == "" {
					fieldDesc = "-"
				}

				// 转义特殊字符防止表格破坏
				safeFieldName := strings.ReplaceAll(fieldName, "|", "\\|")
				safeFieldType := strings.ReplaceAll(fieldType, "|", "\\|")
				// 清理换行符和转义特殊字符
				cleanFieldDesc := strings.ReplaceAll(fieldDesc, "\n", "")
				cleanFieldDesc = strings.ReplaceAll(cleanFieldDesc, "\r", "")
				safeFieldDesc := strings.ReplaceAll(cleanFieldDesc, "|", "\\|")

				md.WriteString(fmt.Sprintf("| %s | %s | %s |\n", safeFieldName, safeFieldType, safeFieldDesc))

				// 如果是对象类型，生成二级表格
				if fieldType == "object" {
					objectProps := value.Get("properties")
					if objectProps.Exists() {
						md.WriteString("\n**" + fieldName + " 对象属性：**\n\n")
						md.WriteString("| 属性名 | 类型 | 说明 |\n")
						md.WriteString("|--------|------|------|\n")
						objectProps.ForEach(func(subKey, subValue gjson.Result) bool {
							propName := subKey.String()
							propType := subValue.Get("type").String()
							propDesc := subValue.Get("description").String()
							if propDesc == "" {
								propDesc = "-"
							}
							// 转义特殊字符防止表格破坏
							safePropName := strings.ReplaceAll(propName, "|", "\\|")
							safePropType := strings.ReplaceAll(propType, "|", "\\|")
							// 清理换行符和转义特殊字符
							cleanPropDesc := strings.ReplaceAll(propDesc, "\n", "")
							cleanPropDesc = strings.ReplaceAll(cleanPropDesc, "\r", "")
							safePropDesc := strings.ReplaceAll(cleanPropDesc, "|", "\\|")

							md.WriteString(fmt.Sprintf("| %s | %s | %s |\n", safePropName, safePropType, safePropDesc))
							return true
						})
						md.WriteString("\n")
					}
				}
				return true
			})
			md.WriteString("\n")

			// 5. 返回值示例
			md.WriteString("## 返回值示例\n\n")
			md.WriteString("**成功响应示例：**\n\n")
			md.WriteString("```json\n")
			md.WriteString("{\n")

			isFirst := true
			endProperties.ForEach(func(key, value gjson.Result) bool {
				fieldName := key.String()

				if !isFirst {
					md.WriteString(",\n")
				}
				isFirst = false

				example := generateExampleValue(value)
				md.WriteString(fmt.Sprintf("  \"%s\": %s", fieldName, example))
				return true
			})

			md.WriteString("\n}\n```\n")
		}
	}

	return md.String()
}

// 原来的生成 Markdown 文档函数（保留作为备用）
func generateMarkdown(nodeData gjson.Result) string {
	var md strings.Builder

	// 获取基本信息
	nodeID := nodeData.Get("id").String()
	nodeType := nodeData.Get("type").String()
	nodeTitle := nodeData.Get("data.title").String()
	positionX := nodeData.Get("meta.position.x").Int()
	positionY := nodeData.Get("meta.position.y").Int()

	// 文档标题
	md.WriteString(fmt.Sprintf("# %s 组件参数文档\n\n", nodeTitle))

	// 组件基本信息
	md.WriteString("## 基本信息\n\n")
	md.WriteString(fmt.Sprintf("- **组件ID**: `%s`\n", nodeID))
	md.WriteString(fmt.Sprintf("- **组件类型**: `%s`\n", nodeType))
	md.WriteString(fmt.Sprintf("- **组件标题**: %s\n", nodeTitle))
	md.WriteString(fmt.Sprintf("- **位置**: x=%d, y=%d\n\n", positionX, positionY))

	// 获取参数信息
	outputs := nodeData.Get("data.outputs")
	properties := outputs.Get("properties")
	required := outputs.Get("required")

	// 统计参数数量
	propertiesCount := 0
	properties.ForEach(func(key, value gjson.Result) bool {
		propertiesCount++
		return true
	})

	requiredCount := 0
	required.ForEach(func(key, value gjson.Result) bool {
		requiredCount++
		return true
	})

	// 参数概览
	md.WriteString("## 参数概览\n\n")
	md.WriteString(fmt.Sprintf("该组件共包含 **%d** 个输出参数，其中 **%d** 个为必需参数。\n\n",
		propertiesCount, requiredCount))

	// 必需参数列表
	if requiredCount > 0 {
		md.WriteString("### 必需参数\n\n")
		required.ForEach(func(key, value gjson.Result) bool {
			paramName := value.String()
			paramData := properties.Get(paramName)
			paramType := paramData.Get("type").String()
			paramDesc := paramData.Get("description").String()
			md.WriteString(fmt.Sprintf("- `%s` (%s): %s\n", paramName, paramType, paramDesc))
			return true
		})
		md.WriteString("\n")
	}

	// 详细参数说明
	md.WriteString("## 详细参数说明\n\n")

	// 收集并排序参数
	type ParamInfo struct {
		Name  string
		Data  gjson.Result
		Index int
	}

	var params []ParamInfo
	properties.ForEach(func(key, value gjson.Result) bool {
		paramName := key.String()
		index := value.Get("extra.index").Int()
		params = append(params, ParamInfo{
			Name:  paramName,
			Data:  value,
			Index: int(index),
		})
		return true
	})

	// 按 index 排序
	sort.Slice(params, func(i, j int) bool {
		return params[i].Index < params[j].Index
	})

	// 生成每个参数的详细说明
	for _, param := range params {
		name := param.Name
		prop := param.Data

		md.WriteString(fmt.Sprintf("### %s\n\n", name))
		md.WriteString(fmt.Sprintf("- **参数名**: `%s`\n", name))
		md.WriteString(fmt.Sprintf("- **类型**: `%s`\n", prop.Get("type").String()))
		md.WriteString(fmt.Sprintf("- **描述**: %s\n", prop.Get("description").String()))

		// 是否必需
		isRequired := isRequiredParam(required, name)
		md.WriteString(fmt.Sprintf("- **是否必需**: %s\n", map[bool]string{true: "是", false: "否"}[isRequired]))

		// 默认值
		if prop.Get("default").Exists() {
			defaultValue := prop.Get("default")
			md.WriteString(fmt.Sprintf("- **默认值**: `%v`\n", defaultValue.Value()))
		}

		// 索引
		if prop.Get("extra.index").Exists() {
			index := prop.Get("extra.index").Int()
			md.WriteString(fmt.Sprintf("- **参数索引**: %d\n", index))
		}

		// 特殊类型处理
		paramType := prop.Get("type").String()
		switch paramType {
		case "object":
			objectProps := prop.Get("properties")
			if objectProps.Exists() {
				md.WriteString("- **对象属性**:\n")
				objectProps.ForEach(func(key, value gjson.Result) bool {
					propName := key.String()
					propType := value.Get("type").String()
					propDesc := value.Get("description").String()
					md.WriteString(fmt.Sprintf("  - `%s` (%s): %s\n", propName, propType, propDesc))
					return true
				})
			}
		case "array":
			itemType := prop.Get("items.type").String()
			if itemType != "" {
				md.WriteString(fmt.Sprintf("- **数组元素类型**: `%s`\n", itemType))
			}
		}

		md.WriteString("\n")
	}

	// 使用示例
	md.WriteString("## 使用示例\n\n")
	md.WriteString("```json\n")
	md.WriteString("{\n")

	for i, param := range params {
		name := param.Name
		prop := param.Data

		if i > 0 {
			md.WriteString(",\n")
		}

		example := generateExampleValue(prop)
		md.WriteString(fmt.Sprintf("  \"%s\": %s", name, example))
	}

	md.WriteString("\n}\n```\n\n")

	// 参数验证规则
	md.WriteString("## 参数验证规则\n\n")
	md.WriteString(fmt.Sprintf("- 必需参数数量: %d\n", requiredCount))
	md.WriteString(fmt.Sprintf("- 总参数数量: %d\n", propertiesCount))
	md.WriteString("- 所有必需参数必须提供值\n")
	md.WriteString("- 参数类型必须严格匹配定义\n\n")

	// 生成时间戳
	md.WriteString("---\n")
	md.WriteString("*文档自动生成于 DSL 解析*\n")

	return md.String()
}

// 检查参数是否为必需参数
func isRequiredParam(required gjson.Result, paramName string) bool {
	isReq := false
	required.ForEach(func(key, value gjson.Result) bool {
		if value.String() == paramName {
			isReq = true
			return false // 找到后停止遍历
		}
		return true
	})
	return isReq
}

// 生成示例值
func generateExampleValue(prop gjson.Result) string {
	paramType := prop.Get("type").String()

	switch paramType {
	case "string":
		if prop.Get("default").Exists() {
			return fmt.Sprintf("\"%s\"", prop.Get("default").String())
		}
		paramName := prop.Get("name").String()
		return fmt.Sprintf("\"%s示例\"", paramName)
	case "boolean":
		if prop.Get("default").Exists() {
			return strconv.FormatBool(prop.Get("default").Bool())
		}
		return "true"
	case "integer", "number":
		if prop.Get("default").Exists() {
			return fmt.Sprintf("%v", prop.Get("default").Value())
		}
		return "0"
	case "array":
		itemType := prop.Get("items.type").String()
		if itemType == "string" {
			return "[\"示例1\", \"示例2\"]"
		}
		return "[]"
	case "object":
		objectProps := prop.Get("properties")
		if objectProps.Exists() {
			var objParts []string
			objectProps.ForEach(func(key, value gjson.Result) bool {
				propName := key.String()
				subExample := generateExampleValue(value)
				objParts = append(objParts, fmt.Sprintf("\"%s\": %s", propName, subExample))
				return true
			})
			return fmt.Sprintf("{ %s }", strings.Join(objParts, ", "))
		}
		return "{}"
	default:
		return "null"
	}
}
