# Etapa 1: Construir o binário
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod download && go mod verify && go mod tidy
RUN go build -v -o main ./main.go

# Etapa 2: Imagem final enxuta
FROM alpine:latest

WORKDIR /app

# Copiar binário da etapa de build
COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

EXPOSE 8080
ENV GIN_MODE=release

CMD ["/app/main"]