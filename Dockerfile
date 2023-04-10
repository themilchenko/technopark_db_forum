# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /app
COPY / .

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN go mod tidy
RUN go build -o main cmd/main.go

CMD ["./main"]
