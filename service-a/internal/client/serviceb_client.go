package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/guilhermeayusso/clima-otel/service-a/internal/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ServiceBClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewServiceBClient(baseURL string) *ServiceBClient {
	return &ServiceBClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *ServiceBClient) GetWeatherByCEP(ctx context.Context, cep string) ([]byte, int, error) {
	payload := dto.WeatherRequest{
		CEP: cep,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}

	url := fmt.Sprintf("%s/weather", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return responseBody, resp.StatusCode, nil
}
