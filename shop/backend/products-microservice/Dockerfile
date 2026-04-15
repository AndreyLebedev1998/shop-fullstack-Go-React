FROM golang:1.26.1-alpine

WORKDIR /products-microservice

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["sh", "-c", "air -c .air.toml"]
