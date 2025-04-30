# Etapa 1: build da aplicação
FROM golang:1.24 AS builder
WORKDIR /app

# Copia os arquivos
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build do binário para Linux amd64 (e saída direta para /server)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server /app/cmd/server/main.go

# Etapa 2: imagem final e enxuta
FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copia o binário
COPY --from=builder /server /app/server

# Copia .env se quiser (opcional)
COPY --from=builder /app/cmd/server/.env /app/.env

# Expõe a porta 8080
EXPOSE 8080

# Comando para rodar o server
CMD ["./server"]
