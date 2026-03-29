# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar dependências primeiro (cache eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código
COPY . .

# Compilar o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# ---- Stage 2: Runtime ----
FROM alpine:3.21

WORKDIR /app

# Certificados SSL (necessário para chamadas HTTPS)
RUN apk --no-cache add ca-certificates

# Copiar apenas o binário compilado
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]