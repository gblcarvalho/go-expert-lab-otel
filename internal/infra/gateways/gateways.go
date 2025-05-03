package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gblcarvalho/go-expert-lab-otel/internal/gateways"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ViaCEPGateway struct {
	tracer trace.Tracer
}

const ViaCEPURL = "https://viacep.com.br/ws/%s/json/"
const WeatherAPIURL = "https://api.weatherapi.com/v1/current.json?q=%s&key=%s"
const GetWeatherURL = "%s/weather-apis/%s"

type ViaCEPResp struct {
	Localidade string `json:"localidade,omitempty"`
	Erro       string `json:"erro,omitempty"`
}

type WeatherAPIResp struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type GetWeatherResp struct {
	City       string  `json:"city,omitempty"`
	Celsius    float64 `json:"temp_C,omitempty"`
	Fahrenheit float64 `json:"temp_F,omitempty"`
	Kelvin     float64 `json:"temp_K,omitempty"`
}

func NewViaCEPGateway(tracer trace.Tracer) *ViaCEPGateway {
	return &ViaCEPGateway{tracer: tracer}
}

func (g *ViaCEPGateway) GetLocation(ctx context.Context, cep string) (gateways.CEPLocation, error) {
	url := fmt.Sprintf(ViaCEPURL, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gateways.CEPLocation{}, fmt.Errorf("erro ao montar requisição: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gateways.CEPLocation{}, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gateways.CEPLocation{}, fmt.Errorf("requisição falhou com status: %s", resp.Status)
	}

	var result ViaCEPResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return gateways.CEPLocation{}, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if result.Erro != "" {
		return gateways.CEPLocation{}, fmt.Errorf("requisição falhou")
	}

	return gateways.CEPLocation{
		Locality: result.Localidade,
	}, nil
}

type WeatherAPIGateway struct {
	apiKey string
	tracer trace.Tracer
}

func NewWeatherAPIGateway(apiKey string, tracer trace.Tracer) *WeatherAPIGateway {
	return &WeatherAPIGateway{apiKey: apiKey, tracer: tracer}
}

func (w *WeatherAPIGateway)	GetWeather(ctx context.Context, locality string) (gateways.WeatherTemp, error) {
	escapedLocality := url.QueryEscape(locality)
	url := fmt.Sprintf(WeatherAPIURL, escapedLocality, w.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gateways.WeatherTemp{}, fmt.Errorf("erro ao montar requisição: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gateways.WeatherTemp{}, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gateways.WeatherTemp{}, fmt.Errorf("requisição falhou com status: %s", resp.Status)
	}

	var result WeatherAPIResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return gateways.WeatherTemp{}, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return gateways.WeatherTemp{
		City: result.Location.Name,
		Celsius: result.Current.TempC,
	}, nil
}

type GetWeatherGateway struct {
	host string
	tracer trace.Tracer
}

func NewGetWeatherGateway(host string, tracer trace.Tracer) *GetWeatherGateway {
	return &GetWeatherGateway{host: host, tracer: tracer}
}

func (w *GetWeatherGateway)	GetWeather(ctx context.Context, cep string) (gateways.GetWeatherResult, error) {
	url := fmt.Sprintf(GetWeatherURL, w.host, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gateways.GetWeatherResult{}, fmt.Errorf("erro ao montar requisição: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gateways.GetWeatherResult{}, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return gateways.GetWeatherResult{}, utils.ErrCEPNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return gateways.GetWeatherResult{}, fmt.Errorf("requisição falhou com status: %s", resp.Status)
	}

	var result GetWeatherResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return gateways.GetWeatherResult{}, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return gateways.GetWeatherResult{
		City: result.City,
		Celsius: result.Celsius,
		Fahrenheit: result.Fahrenheit,
		Kelvin: result.Kelvin,
	}, nil
}
