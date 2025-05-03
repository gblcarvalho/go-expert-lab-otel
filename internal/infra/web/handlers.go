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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WeatherHandler struct {
	cepGateway        gateways.CEPGatewayInterface
	weatherGateway    gateways.WeatherGatewayInterface
	getWeatherGateway gateways.GetWeatherInterface
	tracer            trace.Tracer
}

func NewWeatherHandler(
	cepGateway gateways.CEPGatewayInterface,
	weatherGateway gateways.WeatherGatewayInterface,
	getWeatherGateway gateways.GetWeatherInterface,
	tracer trace.Tracer,
) *WeatherHandler {
	return &WeatherHandler{
		cepGateway:        cepGateway,
		weatherGateway:    weatherGateway,
		getWeatherGateway: getWeatherGateway,
		tracer:            tracer,
	}
}

type GetCompleteReq struct {
	CEP string `json:"cep"`
}

func (h *WeatherHandler) GetComplete(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier) 
	ctx, span := h.tracer.Start(ctx, "get weather complete")
	defer span.End()

	var req GetCompleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.CEP == "" {
		http.Error(w, "cep is required", http.StatusBadRequest)
		return
	}

	getWeatherCompleteUC := usecases.NewGetWeatherCompleteUseCase(h.getWeatherGateway)
	output, err := getWeatherCompleteUC.Execute(ctx, req.CEP)
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
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier) 
	ctx, span := h.tracer.Start(ctx, "get weather")
	defer span.End()

	cep := chi.URLParam(r, "cep")
	getWeatherUC := usecases.NewGetWeatherUseCase(h.cepGateway, h.weatherGateway)
	output, err := getWeatherUC.Execute(ctx, cep)
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
