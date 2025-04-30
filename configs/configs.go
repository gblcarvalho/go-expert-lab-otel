package configs

import "github.com/spf13/viper"

type conf struct {
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
	GetWeatherHost string `mapstructure:"GET_WEATHER_HOST"`
	OTELServiceName string `mapstructure:"OTEL_SERVICE_NAME"`
	OTELCollectorURL string `mapstructure:"OTEL_COLLECTOR_URL"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
