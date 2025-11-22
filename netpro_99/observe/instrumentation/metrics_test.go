package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func setupTestHandler(t *testing.T) http.Handler {
	ctx := context.Background()

	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter("hippo-an")

	counter, _ := meter.Int64Counter("demo_api_requests_total")
	histogram, _ := meter.Float64Histogram("demo_api_duration_seconds")

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, OpenTelemetry!"))

		attrs := metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("route", "/hello"),
		)
		counter.Add(ctx, 1, attrs)
		histogram.Record(ctx, 0.05, attrs)
	})

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func TestHelloEndpoint(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != "Hello, OpenTelemetry!" {
		t.Errorf("expected 'Hello, OpenTelemetry!', got '%s'", string(body))
	}
}

func TestMetricsEndpoint(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/metrics", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	metricsOutput := string(body)

	if !strings.Contains(metricsOutput, "demo_api_requests_total") {
		t.Error("metrics should contain demo_api_requests_total")
	}

	if !strings.Contains(metricsOutput, "demo_api_duration_seconds") {
		t.Error("metrics should contain demo_api_duration_seconds")
	}

	lines := strings.Split(metricsOutput, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "demo_api_requests_total") && strings.Contains(line, `method="GET"`) {
			t.Logf("%s\n", line)
		}
		if strings.HasPrefix(line, "demo_api_duration_seconds") && strings.Contains(line, `method="GET"`) {
			t.Logf("%s\n", line)
		}
	}
}
