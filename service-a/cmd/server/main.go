package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/guilhermeayusso/clima-otel/service-a/internal/client"
	"github.com/guilhermeayusso/clima-otel/service-a/internal/handler"
	"github.com/guilhermeayusso/clima-otel/service-a/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	otelCollectorURL := os.Getenv("OTEL_COLLECTOR_URL")
	if otelCollectorURL == "" {
		otelCollectorURL = "localhost:4318"
	}

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://localhost:8081"
	}

	tp, err := otel.InitTracerProvider(context.Background(), "service-a", otelCollectorURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	serviceBClient := client.NewServiceBClient(serviceBURL)
	weatherHandler := handler.NewWeatherHandler(serviceBClient)

	mux := http.NewServeMux()
	mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(weatherHandler.Handle), "service-a-weather"))

	addr := ":8080"
	log.Println("service-a running on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
