# ---------- Builder ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener .

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]