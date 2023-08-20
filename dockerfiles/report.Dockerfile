FROM golang:1.21

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o reportingapicmd ./cmd/reportingapi

EXPOSE 3000
CMD ["./reportingapicmd"]