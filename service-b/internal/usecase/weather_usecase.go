package usecase

import (
	"context"
	"errors"
	"math"
	"unicode"

	"github.com/guilhermeayusso/clima-otel/service-b/internal/client"
	"github.com/guilhermeayusso/clima-otel/service-b/internal/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var ErrInvalidZipcode = errors.New("invalid zipcode")
var ErrZipcodeNotFound = errors.New("can not find zipcode")

type WeatherUseCase struct {
	viaCEPClient     *client.ViaCEPClient
	weatherAPIClient *client.WeatherAPIClient
}

func NewWeatherUseCase(
	viaCEPClient *client.ViaCEPClient,
	weatherAPIClient *client.WeatherAPIClient,
) *WeatherUseCase {
	return &WeatherUseCase{
		viaCEPClient:     viaCEPClient,
		weatherAPIClient: weatherAPIClient,
	}
}

func (uc *WeatherUseCase) Execute(ctx context.Context, cep string) (dto.WeatherResponse, error) {
	tracer := otel.Tracer("service-b-usecase")
	ctx, span := tracer.Start(ctx, "weather.execute")
	defer span.End()

	span.SetAttributes(attribute.String("zipcode", cep))

	if !isValidCEP(cep) {
		span.RecordError(ErrInvalidZipcode)
		return dto.WeatherResponse{}, ErrInvalidZipcode
	}

	city, err := uc.viaCEPClient.FindCityByCEP(ctx, cep)
	if err != nil {
		if errors.Is(err, client.ErrZipcodeNotFound) {
			span.RecordError(ErrZipcodeNotFound)
			return dto.WeatherResponse{}, ErrZipcodeNotFound
		}
		span.RecordError(err)
		return dto.WeatherResponse{}, err
	}

	span.SetAttributes(attribute.String("city", city))

	tempC, err := uc.weatherAPIClient.GetTemperatureByCity(ctx, city)
	if err != nil {
		span.RecordError(err)
		return dto.WeatherResponse{}, err
	}

	response := dto.WeatherResponse{
		City:  city,
		TempC: roundToOneDecimal(tempC),
		TempF: roundToOneDecimal(celsiusToFahrenheit(tempC)),
		TempK: roundToOneDecimal(celsiusToKelvin(tempC)),
	}

	span.SetAttributes(
		attribute.Float64("temp_c", response.TempC),
		attribute.Float64("temp_f", response.TempF),
		attribute.Float64("temp_k", response.TempK),
	)

	return response, nil
}

func isValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	for _, r := range cep {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func celsiusToFahrenheit(tempC float64) float64 {
	return tempC*1.8 + 32
}

func celsiusToKelvin(tempC float64) float64 {
	return tempC + 273
}

func roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}
