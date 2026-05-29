BINARY   := ./bin/app
MAIN     := ./cmd/main.go
MIGRATE  := migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

.PHONY: run build test tidy docker-up docker-down migrate-up migrate-down migrate-create swag

# Jalankan server dengan Air (live reload)
run:
	air

# Build binary ke ./bin/app
build:
	go build -ldflags="-s -w" -o $(BINARY) $(MAIN)

# Jalankan semua test dengan verbose dan coverage
test:
	go test ./... -v -cover

# Bersihkan dependencies di go.mod dan go.sum
tidy:
	go mod tidy

# Jalankan semua container Docker
docker-up:
	docker compose up -d

# Hentikan semua container Docker
docker-down:
	docker compose down

# Jalankan semua migration yang belum diterapkan
migrate-up:
	$(MIGRATE) up

# Batalkan 1 langkah migration terakhir
migrate-down:
	$(MIGRATE) down 1

# Buat file migration baru: make migrate-create name=nama_migration
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Generate dokumentasi Swagger
swag:
	swag init -g $(MAIN) -o ./docs
