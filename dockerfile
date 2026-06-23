# --- Stage 1: build ---
FROM golang:1.26 AS builder
WORKDIR /app

# Copiar go.mod/go.sum primero y bajar deps (capa cacheable)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto y compilar ambos binarios
COPY . .
RUN go build -o /bin/api ./cmd/api
RUN go build -o /bin/worker ./cmd/worker

# --- Stage 2: runtime ---
FROM debian:bookworm-slim
WORKDIR /app

# Copiar solo los binarios desde el builder
COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/worker /bin/worker

# Por defecto corre el api (el compose sobreescribe para el worker)
CMD ["/bin/api"]