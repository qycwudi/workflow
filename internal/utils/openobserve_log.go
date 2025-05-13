// utils/logger.go
package utils

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"workflow/internal/config"
)

func InitLogger(config config.OpenObserveConfig) *sdklog.LoggerProvider {
	exporter, err := otlploggrpc.New(
		context.Background(),
		otlploggrpc.WithEndpoint(config.OPEN_OBSERVE_ENDPOINT),
		otlploggrpc.WithHeaders(map[string]string{
			"Authorization": config.OPEN_OBSERVE_TOKEN,
			"organization":  config.OPEN_OBSERVE_ORGANIZATION,
			"stream-name":   config.OPEN_OBSERVE_HOST_NAME + "_LOG",
		}),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println(err)
	}

	processor := sdklog.NewBatchProcessor(exporter)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("workflow"),
	)

	return sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)
}

func GetLogger(provider *sdklog.LoggerProvider) log.Logger {
	return provider.Logger(
		"workflow-logger",
		log.WithInstrumentationVersion("v1.0.0"),
	)
}
