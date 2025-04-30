package usecases

import (
	"context"

	"github.com/gblcarvalho/go-expert-lab-otel/internal/gateways"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/utils"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/valueobjects"
)


type GetWeatherCompleteUseCase struct {
	getWeatherGateway gateways.GetWeatherInterface
}

type GetWeatherCompleteUseCaseOutput struct {
	City       string  `json:"city,omitempty"`
	Celsius    float64 `json:"temp_C,omitempty"`
	Fahrenheit float64 `json:"temp_F,omitempty"`
	Kelvin     float64 `json:"temp_K,omitempty"`
}

func NewGetWeatherCompleteUseCase(getWeatherGateway gateways.GetWeatherInterface) *GetWeatherCompleteUseCase {
	return &GetWeatherCompleteUseCase{getWeatherGateway: getWeatherGateway}
}

func (uc *GetWeatherCompleteUseCase) Execute(ctx context.Context, cepStr string) (GetWeatherCompleteUseCaseOutput, error) {
	cep, err := valueobjects.NewCEP(cepStr)
	if err != nil {
		return GetWeatherCompleteUseCaseOutput{}, utils.ErrInvalidCEP
	}
	result, err := uc.getWeatherGateway.GetWeather(ctx, cep.Value())
	if err != nil {
		return GetWeatherCompleteUseCaseOutput{}, err
	}
	return GetWeatherCompleteUseCaseOutput{
		City: result.City,
		Celsius: result.Celsius,
		Fahrenheit: result.Fahrenheit,
		Kelvin: result.Kelvin,
	}, nil
}

type GetWeatherUseCase struct {
	cepGateway     gateways.CEPGatewayInterface
	weatherGateway gateways.WeatherGatewayInterface
}

type GetWeatherUseCaseOutput struct {
	City       string  `json:"city,omitempty"`
	Celsius    float64 `json:"temp_C,omitempty"`
	Fahrenheit float64 `json:"temp_F,omitempty"`
	Kelvin     float64 `json:"temp_K,omitempty"`
}

func NewGetWeatherUseCase(
	cepGateway gateways.CEPGatewayInterface,
	weatherGateway gateways.WeatherGatewayInterface,
) *GetWeatherUseCase {
	return &GetWeatherUseCase{
		cepGateway:     cepGateway,
		weatherGateway: weatherGateway,
	}
}

func (uc *GetWeatherUseCase) Execute(ctx context.Context, cepStr string) (GetWeatherUseCaseOutput, error) {
	location, err := uc.cepGateway.GetLocation(ctx, cepStr)
	if err != nil {
		return GetWeatherUseCaseOutput{}, utils.ErrCEPNotFound
	}
	weather, err := uc.weatherGateway.GetWeather(ctx, location.Locality)
	if err != nil {
		return GetWeatherUseCaseOutput{}, utils.ErrWeather
	}

	return GetWeatherUseCaseOutput{
		Celsius: weather.Celsius,
		Fahrenheit: CelsiusToFahrenheit(weather.Celsius),
		Kelvin: CelsiusToKelvin(weather.Celsius),
	}, nil
}

func CelsiusToFahrenheit(tempC float64) float64 {
	return tempC * 1.8 + 32
}

func CelsiusToKelvin(tempC float64) float64 {
	return tempC + 273
}
