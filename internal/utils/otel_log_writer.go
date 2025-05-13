package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/log"
)

type OtelLogWriter struct {
	logger log.Logger
}

func NewOtelLogWriter(logger log.Logger) *OtelLogWriter {
	return &OtelLogWriter{
		logger: logger, // 使用小写字段但通过构造函数设置
	}
}

// 必须实现所有接口方法
func (w *OtelLogWriter) Alert(v any) {
	w.emit(log.SeverityError, fmt.Sprint(v))
}

func (w *OtelLogWriter) Close() error {
	return nil
}

func (w *OtelLogWriter) Debug(v any, fields ...logx.LogField) {
	w.emitWithFields(log.SeverityDebug, fmt.Sprint(v), fields)
}

func (w *OtelLogWriter) Error(v any, fields ...logx.LogField) {
	w.emitWithFields(log.SeverityError, fmt.Sprint(v), fields)
}

func (w *OtelLogWriter) Info(v any, fields ...logx.LogField) {
	w.emitWithFields(log.SeverityInfo, fmt.Sprint(v), fields)
}

func (w *OtelLogWriter) Severe(v any) {
	w.emit(log.SeverityFatal, fmt.Sprint(v))
}

func (w *OtelLogWriter) Slow(v any, fields ...logx.LogField) {
	w.emitWithFields(log.SeverityTrace, fmt.Sprint(v), fields)
}

func (w *OtelLogWriter) Stack(v any) {
	w.emitWithFields(log.SeverityError, fmt.Sprint(v),
		[]logx.LogField{
			{Key: "stack", Value: "true"},
		},
	)
}

func (w *OtelLogWriter) Stat(v any, fields ...logx.LogField) {
	w.emitWithFields(log.SeverityInfo, fmt.Sprint(v), fields)
}

// 基础emit方法
func (w *OtelLogWriter) emit(severity log.Severity, message string) {
	rec := log.Record{}
	rec.SetTimestamp(time.Now())
	rec.SetObservedTimestamp(time.Now())
	rec.SetSeverity(severity)
	rec.SetBody(log.StringValue(message))
	w.logger.Emit(context.Background(), rec)
}

// 带字段的emit方法
func (w *OtelLogWriter) emitWithFields(severity log.Severity, message string, fields []logx.LogField) {
	rec := log.Record{}
	rec.SetTimestamp(time.Now())
	rec.SetObservedTimestamp(time.Now())
	rec.SetSeverity(severity)
	rec.SetBody(log.StringValue(message))

	attrs := make([]log.KeyValue, 0, len(fields))
	for _, f := range fields {
		attrs = append(attrs, log.String(f.Key, fmt.Sprint(f.Value)))
	}
	rec.AddAttributes(attrs...)

	w.logger.Emit(context.Background(), rec)
}

// 实现Write接口
func (w *OtelLogWriter) Write(p []byte) (n int, err error) {
	w.emit(log.SeverityInfo, string(p))
	return len(p), nil
}
