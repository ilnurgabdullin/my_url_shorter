FROM golang:1.20

WORKDIR /app

ENV TZ=Europe/Moscow
ENV DB_NAME=GO_DATA
ENV DB_HOST=db
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=123

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
