# syntax=docker/dockerfile:1
FROM golang:1.24.3-alpine3.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]