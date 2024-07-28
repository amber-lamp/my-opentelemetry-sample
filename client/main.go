package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			"example-client",
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to perform request: %v", err)
	}
	defer res.Body.Close()

	log.Println("Response received:", res.Status)
}
