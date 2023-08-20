FROM golang:1.21

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o processorcmd ./cmd/messageprocessor

CMD ["./processorcmd"]