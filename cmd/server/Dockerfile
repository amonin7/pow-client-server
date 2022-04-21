FROM golang:1.17

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o server ./cmd/server

CMD ["/build/server"]