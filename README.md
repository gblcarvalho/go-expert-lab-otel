# go-expert-lab-otel


Implementação de um desafio de OTEL
🛠️ Pré-requisitos

Antes de executar o projeto, certifique-se de ter os seguintes softwares instalados:

- Docker
- Docker Compose

## 🚀 Como rodar o projeto
1. Subindo o ambiente com Docker Compose

Execute os seguintes comando apra subir o ambiente via docker compose.

Primeiro é necessário criar um arquivo .env contendo as variáveis de ambiente. Ele pode ser criado copiando o .env.exemple e alterando o valor de WEATHER_API_KEY

```bash
cp ./cmd/server/.env.exemple ./cmd/server/.env
```

Depois é só executar o comando abaixo passando o arquivo .env criado como parâmetro

```bash
docker compose --env-file ./cmd/server/.env up -d
```

Este comando irá:

- Iniciar dois containers do servidor da applicação Go
- Iniciar OTEL Collector
- Zipkin


O serviço da aplicação está disponíveis em
- HTTP (RESTful): http://localhost:8080/weather-complete/{cep} // Serviço A
- HTTP (RESTful): http://localhost:8080/weather/{cep} // Serviço B

Para acessar o Zipkin
- http://localhost:9411
