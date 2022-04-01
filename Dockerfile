FROM golang:1.17.8-alpine AS builder

ADD . /app
WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /cinetickets ./cmd/cinetickets/main.go

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
COPY --from=builder /cinetickets ./
RUN chmod +x ./cinetickets
ENTRYPOINT ["./cinetickets"]