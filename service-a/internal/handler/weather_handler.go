package handler

import (
	"encoding/json"
	"net/http"

	"github.com/guilhermeayusso/clima-otel/service-a/internal/client"
	"github.com/guilhermeayusso/clima-otel/service-a/internal/dto"
)

type WeatherHandler struct {
	serviceBClient *client.ServiceBClient
}

func NewWeatherHandler(serviceBClient *client.ServiceBClient) *WeatherHandler {
	return &WeatherHandler{
		serviceBClient: serviceBClient,
	}
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request dto.WeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	responseBody, statusCode, err := h.serviceBClient.GetWeatherByCEP(r.Context(), request.CEP)
	if err != nil {
		http.Error(w, "failed to call service-b", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBody)
}
