FROM golang:1.17.8-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o cinetickets ./cmd/cinetickets/main.go

ENTRYPOINT ["/app/cinetickets"]