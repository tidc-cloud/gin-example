# Use the official Golang image to create a build artifact.
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o docker-gs-ping


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/docker-gs-ping ./docker-gs-ping

EXPOSE 8080

CMD ["/app/docker-gs-ping"]

