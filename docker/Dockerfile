FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /quoter cmd/quoter/main.go

EXPOSE 8080

CMD ["/quoter"]