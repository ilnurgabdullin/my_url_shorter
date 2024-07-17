# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
COPY . .

RUN go build -o main .

CMD ["./main"]
