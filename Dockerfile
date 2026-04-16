FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bonsai .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bonsai .
COPY templates/ templates/
COPY static/ static/
EXPOSE 3100
CMD ["./bonsai"]
