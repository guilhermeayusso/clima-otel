package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/guilhermeayusso/clima-otel/service-b/internal/client"
	"github.com/guilhermeayusso/clima-otel/service-b/internal/handler"
	"github.com/guilhermeayusso/clima-otel/service-b/internal/otel"
	"github.com/guilhermeayusso/clima-otel/service-b/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	otelCollectorURL := os.Getenv("OTEL_COLLECTOR_URL")
	if otelCollectorURL == "" {
		otelCollectorURL = "localhost:4318"
	}

	tp, err := otel.InitTracerProvider(context.Background(), "service-b", otelCollectorURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	viaCEPClient := client.NewViaCEPClient()
	weatherAPIClient := client.NewWeatherAPIClient(weatherAPIKey)

	weatherUseCase := usecase.NewWeatherUseCase(viaCEPClient, weatherAPIClient)
	weatherHandler := handler.NewWeatherHandler(weatherUseCase)

	mux := http.NewServeMux()
	mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(weatherHandler.Handle), "service-b-weather"))

	addr := ":8081"
	log.Println("service-b running on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
