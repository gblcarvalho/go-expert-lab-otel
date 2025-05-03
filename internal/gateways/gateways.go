package gateways

import "context"

type GetWeatherInterface interface {
	GetWeather(ctx context.Context, cep string) (GetWeatherResult, error)
}

type CEPGatewayInterface interface {
	GetLocation(ctx context.Context, cep string) (CEPLocation, error)
}

type WeatherGatewayInterface interface {
	GetWeather(ctx context.Context, locality string) (WeatherTemp, error)
}

type CEPLocation struct {
	Locality string
}

type WeatherTemp struct {
	City    string
	Celsius float64
}

type GetWeatherResult struct {
	City       string
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}
