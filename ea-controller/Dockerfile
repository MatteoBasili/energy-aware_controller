FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy 
RUN go build -o ea-controller ./cmd/controller

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ea-controller .
CMD ["./ea-controller"]

