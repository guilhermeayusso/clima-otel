package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewWeatherAPIClient(apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		baseURL: "https://api.weatherapi.com/v1/current.json",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *WeatherAPIClient) GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b-client")
	ctx, span := tracer.Start(ctx, "weatherapi.lookup")
	defer span.End()

	span.SetAttributes(attribute.String("city", city))

	escapedCity := url.QueryEscape(city)
	fullURL := fmt.Sprintf("%s?key=%s&q=%s", c.baseURL, c.apiKey, escapedCity)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	defer resp.Body.Close()

	var response WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		span.RecordError(err)
		return 0, err
	}

	span.SetAttributes(attribute.Float64("temp_c", response.Current.TempC))

	return response.Current.TempC, nil
}
