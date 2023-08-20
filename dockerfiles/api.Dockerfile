FROM golang:1.21

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o apicmd ./cmd/api

EXPOSE 8080
CMD ["./apicmd"]