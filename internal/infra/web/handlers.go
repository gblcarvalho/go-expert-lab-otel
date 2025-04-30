package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gblcarvalho/go-expert-lab-otel/internal/gateways"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/usecases"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/utils"
	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	cepGateway     gateways.CEPGatewayInterface
	weatherGateway gateways.WeatherGatewayInterface
	getWeatherGateway gateways.GetWeatherInterface
}

func NewWeatherHandler(
	cepGateway gateways.CEPGatewayInterface,
	weatherGateway gateways.WeatherGatewayInterface,
	getWeatherGateway gateways.GetWeatherInterface,
) *WeatherHandler {
	return &WeatherHandler{
		cepGateway:     cepGateway,
		weatherGateway: weatherGateway,
		getWeatherGateway: getWeatherGateway,
	}
}

func (h *WeatherHandler) GetComplete(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	getWeatherCompleteUC := usecases.NewGetWeatherCompleteUseCase(h.getWeatherGateway)
	output, err := getWeatherCompleteUC.Execute(cep)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, utils.ErrInvalidCEP) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		} else if errors.Is(err, utils.ErrCEPNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, "can not find weather", http.StatusServiceUnavailable)
		}
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	getWeatherUC := usecases.NewGetWeatherUseCase(h.cepGateway, h.weatherGateway)
	output, err := getWeatherUC.Execute(cep)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCEP) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		} else if errors.Is(err, utils.ErrCEPNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, "can not find weather", http.StatusServiceUnavailable)
		}
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
