# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o ./bin/app ./cmd/main.go

# Stage 2: Image runtime minimal
FROM alpine:latest AS runner

# Instal sertifikat CA untuk koneksi HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin binary dari stage builder
COPY --from=builder /app/bin/app .

# Jalankan dengan user non-root demi keamanan
RUN adduser -D -u 1001 appuser
USER appuser

EXPOSE 8080

CMD ["./app"]
