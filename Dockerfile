# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ratchet-gallery ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/ratchet-gallery .
COPY templates ./templates
COPY static ./static
COPY content ./content

EXPOSE 8070

CMD ["./ratchet-gallery"]
