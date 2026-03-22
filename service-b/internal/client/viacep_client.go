package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

type ViaCEPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewViaCEPClient() *ViaCEPClient {
	return &ViaCEPClient{
		baseURL: "https://viacep.com.br/ws",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ViaCEPClient) FindCityByCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b-client")
	ctx, span := tracer.Start(ctx, "viacep.lookup")
	defer span.End()

	span.SetAttributes(attribute.String("zipcode", cep))

	url := fmt.Sprintf("%s/%s/json/", c.baseURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	var response ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		span.RecordError(err)
		return "", err
	}

	if response.Erro == "true" || response.Localidade == "" {
		span.RecordError(ErrZipcodeNotFound)
		return "", ErrZipcodeNotFound
	}

	span.SetAttributes(attribute.String("city", response.Localidade))

	return response.Localidade, nil
}
