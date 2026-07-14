.PHONY: dev run build test fmt vet swagger docker down db-reset clean

dev:
	air

run:
	go run .

build:
	go build -o url-shortener .

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

swagger:
	swag init

docker:
	docker compose up --build

docker-bg:
	docker compose up -d --build

down:
	docker compose down

db-reset:
	docker compose down -v

clean:
	rm -f url-shortener

check: fmt vet test