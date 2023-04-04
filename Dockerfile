# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

RUN apk update && apk add git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go mod download github.com/golang-jwt/jwt
RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main"]
