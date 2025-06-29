package utils

import (
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogLevel 定义日志级别常量
const (
	ModuleEngine     = "engine"
	ModuleComponent  = "component"
	ModuleDatabase   = "database"
	ModuleAPI        = "api"
	ModuleJob        = "job"
	ModuleAuth       = "auth"
	ModuleWorkflow   = "workflow"
	ModuleCanvas     = "canvas"
	ModuleDatasource = "datasource"
)

// OperationType 定义操作类型常量
const (
	OpCreate   = "create"
	OpUpdate   = "update"
	OpDelete   = "delete"
	OpQuery    = "query"
	OpExecute  = "execute"
	OpPublish  = "publish"
	OpLogin    = "login"
	OpLogout   = "logout"
	OpRegister = "register"
)

// LogBusinessOperation 记录业务操作日志
func LogBusinessOperation(operation, workspaceId, userId string, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("operation", operation),
		logx.Field("workspace_id", workspaceId),
		logx.Field("user_id", userId),
		logx.Field("timestamp", time.Now().Unix()),
	}
	fields = append(fields, extra...)
	logx.Infow("business operation", fields...)
}

// LogAPIRequest 记录API请求日志
func LogAPIRequest(method, path, userId, clientIP string, duration int64, statusCode int) {
	logx.Infow("api request completed",
		logx.Field("method", method),
		logx.Field("path", path),
		logx.Field("user_id", userId),
		logx.Field("client_ip", clientIP),
		logx.Field("duration_ms", duration),
		logx.Field("status_code", statusCode))
}

// LogAPIStart 记录API请求开始
func LogAPIStart(method, path, userId, clientIP string) {
	logx.Infow("api request started",
		logx.Field("method", method),
		logx.Field("path", path),
		logx.Field("user_id", userId),
		logx.Field("client_ip", clientIP))
}

// LogModuleError 记录模块错误日志
func LogModuleError(module, operation string, err error, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("module", module),
		logx.Field("operation", operation),
		logx.Field("error", err.Error()),
	}
	fields = append(fields, extra...)
	logx.Errorw("module operation failed", fields...)
}

// LogWorkflowExecution 记录工作流执行日志
func LogWorkflowExecution(workspaceId, workflowId, traceId, userId string, status string, duration int64, nodeCount int, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("workspace_id", workspaceId),
		logx.Field("workflow_id", workflowId),
		logx.Field("trace_id", traceId),
		logx.Field("user_id", userId),
		logx.Field("status", status),
		logx.Field("duration_ms", duration),
		logx.Field("node_count", nodeCount),
	}
	fields = append(fields, extra...)
	logx.Infow("workflow execution completed", fields...)
}

// LogComponentExecution 记录组件执行日志
func LogComponentExecution(nodeId, nodeType string, status string, duration int64, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("node_id", nodeId),
		logx.Field("node_type", nodeType),
		logx.Field("status", status),
		logx.Field("duration_ms", duration),
	}
	fields = append(fields, extra...)
	logx.Infow("component execution completed", fields...)
}

// LogDatabaseOperation 记录数据库操作日志
func LogDatabaseOperation(operation, table string, duration int64, rowsAffected int64, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("operation", operation),
		logx.Field("table", table),
		logx.Field("duration_ms", duration),
		logx.Field("rows_affected", rowsAffected),
	}
	fields = append(fields, extra...)
	logx.Infow("database operation completed", fields...)
}

// LogAuth 记录认证相关日志
func LogAuth(operation, userId, clientIP, userAgent string, success bool, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("operation", operation),
		logx.Field("user_id", userId),
		logx.Field("client_ip", clientIP),
		logx.Field("user_agent", userAgent),
		logx.Field("success", success),
	}
	fields = append(fields, extra...)
	logx.Infow("auth operation", fields...)
}

// LogPerformanceMetrics 记录性能指标
func LogPerformanceMetrics(module, operation string, duration int64, memoryUsage int64, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("module", module),
		logx.Field("operation", operation),
		logx.Field("duration_ms", duration),
		logx.Field("memory_usage_bytes", memoryUsage),
	}
	fields = append(fields, extra...)
	logx.Infow("performance metrics", fields...)
}

// LogSystemEvent 记录系统事件
func LogSystemEvent(event, component string, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("event", event),
		logx.Field("component", component),
		logx.Field("timestamp", time.Now().Unix()),
	}
	fields = append(fields, extra...)
	logx.Infow("system event", fields...)
}

// LogDebugInfo 记录调试信息（仅在调试模式下）
func LogDebugInfo(module, operation string, data any, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("module", module),
		logx.Field("operation", operation),
		logx.Field("data", data),
	}
	fields = append(fields, extra...)
	logx.Debugw("debug info", fields...)
}

// LogSecurityEvent 记录安全事件
func LogSecurityEvent(event, userId, clientIP, details string, severity string, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("event", event),
		logx.Field("user_id", userId),
		logx.Field("client_ip", clientIP),
		logx.Field("details", details),
		logx.Field("severity", severity),
		logx.Field("timestamp", time.Now().Unix()),
	}
	fields = append(fields, extra...)
	logx.Errorw("security event", fields...)
}

// LogResourceUsage 记录资源使用情况
func LogResourceUsage(resource string, used, total, percentage float64, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("resource", resource),
		logx.Field("used", used),
		logx.Field("total", total),
		logx.Field("percentage", percentage),
	}
	fields = append(fields, extra...)
	
	if percentage > 80 {
		logx.Errorw("high resource usage", fields...)
	} else {
		logx.Infow("resource usage", fields...)
	}
}

// LogJobExecution 记录任务执行日志
func LogJobExecution(jobName, jobType string, status string, duration int64, extra ...logx.LogField) {
	fields := []logx.LogField{
		logx.Field("job_name", jobName),
		logx.Field("job_type", jobType),
		logx.Field("status", status),
		logx.Field("duration_ms", duration),
	}
	fields = append(fields, extra...)
	logx.Infow("job execution completed", fields...)
}