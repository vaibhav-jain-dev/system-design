# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sd-dashboard .

# Runtime stage — minimal Alpine (no Playwright/Chromium; PDF engine is not yet active)
FROM alpine:3.19

# ca-certificates for HTTPS, tzdata for timezone, python3 for future PDF engine
RUN apk add --no-cache ca-certificates tzdata python3

WORKDIR /app
COPY --from=builder /sd-dashboard /app/sd-dashboard

# Copy content + engine (these may also be mounted as volumes)
COPY content/ /app/content/
COPY engine/ /app/engine/

# Output directory for generated PDFs
RUN mkdir -p /app/output

EXPOSE 8080

CMD ["/app/sd-dashboard"]
