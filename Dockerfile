FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

RUN golangci-lint run -v

RUN go test -v -race -cover ./...

RUN go build -o url-shortener ./cmd/url-shortener

CMD ["/app/url-shortener"]