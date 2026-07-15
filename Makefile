.PHONY: dev run build test fmt vet swagger docker down db-reset clean create-admin

dev:
	air

dev-app:
	cd frontend && npm run dev

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
	docker compose -f docker-compose.dev.yml up --build

docker-bg:
	docker compose -f docker-compose.dev.yml up -d --build

down:
	docker compose -f docker-compose.dev.yml down

db-reset:
	docker compose -f docker-compose.dev.yml down -v

clean:
	rm -f url-shortener

create-admin:
	go run scripts/create_admin.go

check: fmt vet test