FROM golang:1.26.1-alpine

WORKDIR /orders-migration-service

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

CMD ["sh", "-c", "goose -dir ./migrations postgres \"$DB_CONN\" up"]
