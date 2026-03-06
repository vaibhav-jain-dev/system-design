.PHONY: run build dev clean

# Run the dashboard
run:
	go run main.go

# Build binary
build:
	go build -o sd-dashboard .

# Run with Docker
dev:
	docker compose up --build

# Clean build artifacts
clean:
	rm -f sd-dashboard
	rm -f output/*.pdf
