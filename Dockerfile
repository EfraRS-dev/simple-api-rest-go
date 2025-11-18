# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY main.go ./

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar el binario compilado
COPY --from=builder /app/main .

# Exponer el puerto
EXPOSE 5001

# Ejecutar la aplicación
CMD ["./main"]
