# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /sd-dashboard .

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache python3 py3-pip chromium

WORKDIR /app
COPY --from=builder /sd-dashboard /app/sd-dashboard

# Install Python deps for PDF engine
COPY engine/requirements.txt /app/engine/requirements.txt
RUN pip3 install --break-system-packages -r /app/engine/requirements.txt && \
    python3 -m playwright install chromium

# Copy content + engine (these may also be mounted as volumes)
COPY content/ /app/content/
COPY engine/ /app/engine/

# Output directory for generated PDFs
RUN mkdir -p /app/output

EXPOSE 8080

CMD ["/app/sd-dashboard"]
