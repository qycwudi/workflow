# 日志和错误处理规范

## 概述

本文档定义了 workflow-back 项目的日志打印和错误处理规范，旨在提高系统的可观测性和维护性。

## 基本原则

### 1. 语言规范
- **日志消息**: 系统内部日志使用英文，便于国际化和运维
- **错误返回**: logic 层返回给前端的错误消息使用中文，提供良好的用户体验
- **代码注释**: 优先使用英文，业务逻辑相关注释可使用中文

### 2. 日志级别使用标准

#### ERROR
- 系统错误，需要立即关注
- 业务关键功能失败
- 第三方服务调用失败
- 数据库连接/操作失败

#### WARN  
- 非关键功能降级
- 资源使用预警
- 配置缺失但有默认值
- 重试机制触发

#### INFO
- 关键业务操作开始/完成
- 状态变更
- 用户重要操作
- 系统启动/关闭

#### DEBUG
- 详细的执行流程
- 参数值和中间结果
- 仅在开发/测试环境使用

## 日志格式规范

### 1. 结构化日志
优先使用结构化日志，便于日志聚合和分析：

```go
// 推荐
logx.Infow("workflow execution started",
    logx.Field("workspace_id", workspaceId),
    logx.Field("workflow_id", workflowId),
    logx.Field("user_id", userId),
    logx.Field("trace_id", traceId))

// 不推荐
logx.Infof("工作流 %s 开始执行，用户: %s", workflowId, userId)
```

### 2. 模块标识
使用统一的模块标识前缀：

```go
// 引擎模块
logx.Infow("engine: workflow compiled successfully", ...)

// 组件模块  
logx.Infow("component: http request completed", ...)

// 数据库模块
logx.Errorw("database: connection failed", ...)

// API层
logx.Infow("api: request received", ...)
```

### 3. 关键字段标准
常用字段名称保持一致：

- `workspace_id`: 工作空间ID
- `user_id`: 用户ID  
- `workflow_id`: 工作流ID
- `node_id`: 节点ID
- `trace_id`: 追踪ID
- `duration_ms`: 执行耗时（毫秒）
- `operation`: 操作类型
- `status`: 状态
- `error`: 错误信息

## 错误处理规范

### 1. Logic 层错误返回
Logic 层返回给前端的错误消息使用中文：

```go
// 参数校验错误
if req.WorkflowId == "" {
    return nil, errors.New(int(logic.ParamError), "工作流ID不能为空")
}

// 业务逻辑错误
if !hasPermission {
    return nil, errors.New(int(logic.AuthError), "您没有权限执行此操作")
}

// 系统错误
if err != nil {
    logx.Errorw("database operation failed", logx.Field("error", err))
    return nil, errors.New(int(logic.SystemError), "系统异常，请稍后重试")
}
```

### 2. 系统内部错误
引擎和组件层的内部错误使用英文：

```go
// 引擎错误
return &EngineExecuteError{Message: "workflow not found: " + workflowID}

// 组件错误  
return nil, errors.New("http component: request timeout")
```

### 3. 错误上下文
记录足够的上下文信息便于问题排查：

```go
if err != nil {
    logx.Errorw("database query failed",
        logx.Field("operation", "find_canvas"),
        logx.Field("workspace_id", workspaceId),
        logx.Field("sql", sql),
        logx.Field("error", err))
    return nil, errors.New(int(logic.SystemStoreError), "查询失败")
}
```

## 性能优化

### 1. 条件日志
避免不必要的字符串格式化：

```go
// 推荐
if logx.DebugEnabled() {
    logx.Debugw("complex debug info", logx.Field("data", complexObject))
}

// 不推荐
logx.Debugf("debug: %+v", complexObject) // 即使debug关闭也会格式化
```

### 2. 敏感信息处理
避免记录敏感信息：

```go
// 不要记录
- 用户密码
- API密钥
- 完整的SQL语句（可能包含敏感数据）
- 用户个人信息

// 可以记录
- 操作类型
- 资源ID
- 执行状态
- 错误类型
```

## 审计日志

关键业务操作需要记录审计日志：

```go
// 用户登录
logx.Infow("user login",
    logx.Field("user_id", userId),
    logx.Field("ip", clientIP),
    logx.Field("user_agent", userAgent))

// 工作流发布
logx.Infow("workflow published",
    logx.Field("workspace_id", workspaceId),
    logx.Field("workflow_id", workflowId),
    logx.Field("user_id", userId),
    logx.Field("version", version))

// 数据源配置变更
logx.Infow("datasource config updated",
    logx.Field("datasource_id", dsId),
    logx.Field("user_id", userId),
    logx.Field("operation", "update_config"))
```

## 日志工具包

创建统一的日志工具包 `pkg/utils/logger.go`：

```go
package utils

import "github.com/zeromicro/go-zero/core/logx"

// 业务操作日志
func LogBusinessOperation(operation, workspaceId, userId string, extra ...logx.Field) {
    fields := []logx.Field{
        logx.Field("operation", operation),
        logx.Field("workspace_id", workspaceId),
        logx.Field("user_id", userId),
    }
    fields = append(fields, extra...)
    logx.Infow("business operation", fields...)
}

// API请求日志
func LogAPIRequest(method, path, userId, clientIP string, duration int64) {
    logx.Infow("api request",
        logx.Field("method", method),
        logx.Field("path", path),
        logx.Field("user_id", userId),
        logx.Field("client_ip", clientIP),
        logx.Field("duration_ms", duration))
}

// 错误日志
func LogError(module, operation string, err error, extra ...logx.Field) {
    fields := []logx.Field{
        logx.Field("module", module),
        logx.Field("operation", operation),
        logx.Field("error", err.Error()),
    }
    fields = append(fields, extra...)
    logx.Errorw("operation failed", fields...)
}
```

## 迁移指南

### 1. 优先级
1. 修改 logic 层错误返回消息为中文
2. 统一引擎和组件层日志为英文
3. 移除过度详细的调试日志
4. 添加关键业务操作审计日志

### 2. 渐进式改进
- 新代码严格遵循此规范
- 现有代码在维护时逐步改进
- 关键模块优先改进

### 3. 工具支持
- 使用IDE插件检查日志格式
- 配置日志分析工具监控日志质量
- 建立日志质量度量指标

## 监控集成

### 1. 指标日志
记录关键业务指标：

```go
logx.Infow("workflow_execution_metrics",
    logx.Field("workspace_id", workspaceId),
    logx.Field("duration_ms", duration),
    logx.Field("node_count", nodeCount),
    logx.Field("status", "success"))
```

### 2. 告警配置
基于日志内容配置告警：
- ERROR级别日志自动告警
- 关键业务操作失败告警
- 性能指标异常告警

这个规范将有助于提高系统的可观测性、问题排查效率和用户体验。