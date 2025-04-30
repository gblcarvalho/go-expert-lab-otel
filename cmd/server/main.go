package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gblcarvalho/go-expert-lab-otel/configs"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/infra/gateways"
	"github.com/gblcarvalho/go-expert-lab-otel/internal/infra/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()
	conn, err := grpc.NewClient(
		collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider.Shutdown, nil
}

func main() {

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	shutdown, err := initProvider(configs.OTELServiceName, configs.OTELCollectorURL)
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)


	tracer := otel.Tracer("go-expert-tracer")
	cepGateway := gateways.NewViaCEPGateway(tracer)
	weatherGateway := gateways.NewWeatherAPIGateway(configs.WeatherApiKey, tracer)
	getWeatherGateway := gateways.NewGetWeatherGateway(configs.GetWeatherHost, tracer)

	handler := web.NewWeatherHandler(cepGateway, weatherGateway, getWeatherGateway, tracer)

	r.Get("/weather-complete/{cep}", handler.GetComplete)
	r.Get("/weather/{cep}", handler.Get)
	r.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":8080", r)
	fmt.Println(err)
}
