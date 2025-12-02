FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN go build -o ecommerce-app ./cmd/main.go

CMD ["./ecommerce-app"]