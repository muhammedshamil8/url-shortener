.PHONY: dev run build test fmt vet swagger docker down db-reset clean

dev:
	air

run:
	go run .

build:
	go build -o url-shortener .

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...

test-coverprofile:
	go tool cover -html=coverage.out && brave-browser http://localhost:8080/coverage/index.html

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