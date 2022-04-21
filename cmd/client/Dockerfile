FROM golang:1.17

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o client ./cmd/client

CMD ["/build/client"]