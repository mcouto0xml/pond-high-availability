# =============================================================================
# Stage 1 — build
# =============================================================================
FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copiar módulos primeiro para aproveitar cache de camadas
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Build usando o funcframework como entrypoint
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o server \
    github.com/GoogleCloudPlatform/functions-framework-go/funcframework

# =============================================================================
# Stage 2 — runtime
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/server /server

# Nome da função que o framework vai registrar
ENV FUNCTION_TARGET=PostTelemetry
ENV PORT=8080

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]