package gateways

type GetWeatherInterface interface {
	GetWeather(cep string) (GetWeatherResult, error)
}

type CEPGatewayInterface interface {
	GetLocation(cep string) (CEPLocation, error)
}

type WeatherGatewayInterface interface {
	GetWeather(locality string) (WeatherTemp, error)
}

type CEPLocation struct {
	Locality string
}

type WeatherTemp struct {
	Celsius    float64
}

type GetWeatherResult struct {
	City       string
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}
