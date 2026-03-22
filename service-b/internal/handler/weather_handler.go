package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guilhermeayusso/clima-otel/service-b/internal/dto"
	"github.com/guilhermeayusso/clima-otel/service-b/internal/usecase"
)

type WeatherHandler struct {
	useCase *usecase.WeatherUseCase
}

func NewWeatherHandler(useCase *usecase.WeatherUseCase) *WeatherHandler {
	return &WeatherHandler{
		useCase: useCase,
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

	response, err := h.useCase.Execute(r.Context(), request.CEP)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidZipcode):
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, usecase.ErrZipcodeNotFound):
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
