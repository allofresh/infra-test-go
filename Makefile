# ============================================================
# Makefile - Go Application Build & Test Targets
# Maps to run_unit_test_components:
#   ["vet", "unit-test", "race-test", "coverage"]
# ============================================================

APP_NAME    := infra-test-go
CMD_PATH    := ./cmd/api
COVER_DIR   := ./coverage
COVER_FILE  := $(COVER_DIR)/coverage.out
COVER_HTML  := $(COVER_DIR)/coverage.html
MIN_COVERAGE := 80

.PHONY: all compile build vet unit-test race-test coverage test clean docker-build docker-run help

## all: Run all test components in pipeline order
all: vet unit-test race-test coverage

OUTPUT_DIR  := deploy/_output/rest
BINARY_NAME := prod

## compile: Download dependencies and build binary
compile:
	@echo "==> Downloading dependencies..."
	go mod download
	@echo "==> Compiling $(APP_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -ldflags="-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "==> Binary: $(OUTPUT_DIR)/$(BINARY_NAME)"
	@echo "==> Compile passed"

## build: Compile the application binary
build:
	@echo "==> Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(CMD_PATH)

## vet: Run go vet for static analysis
vet:
	@echo "==> Running go vet..."
	go vet ./...
	@echo "==> go vet passed"

## unit-test: Run all unit tests
unit-test:
	@echo "==> Running unit tests..."
	go test -v -count=1 ./...
	@echo "==> Unit tests passed"

## race-test: Run tests with race detector enabled
race-test:
	@echo "==> Running tests with race detector..."
	go test -v -race -count=1 ./...
	@echo "==> Race tests passed"

## coverage: Run tests with coverage and enforce minimum threshold
coverage:
	@echo "==> Running tests with coverage..."
	@mkdir -p $(COVER_DIR)
	go test -coverprofile=$(COVER_FILE) -covermode=atomic ./...
	go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@echo "==> Coverage report generated: $(COVER_HTML)"
	@echo "==> Coverage summary:"
	@go tool cover -func=$(COVER_FILE)
	@echo ""
	@TOTAL=$$(go tool cover -func=$(COVER_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "==> Total coverage: $${TOTAL}%"; \
	if [ $$(echo "$${TOTAL} < $(MIN_COVERAGE)" | bc -l) -eq 1 ]; then \
		echo "==> FAIL: Coverage $${TOTAL}% is below minimum $(MIN_COVERAGE)%"; \
		exit 1; \
	else \
		echo "==> PASS: Coverage $${TOTAL}% meets minimum $(MIN_COVERAGE)%"; \
	fi

## test: Alias for running all test components
test: all

## clean: Remove build artifacts and coverage reports
clean:
	@echo "==> Cleaning..."
	rm -rf bin/ $(COVER_DIR)
	@echo "==> Clean complete"

## docker-build: Build the Docker image
docker-build:
	@echo "==> Building Docker image..."
	docker build -t $(APP_NAME):latest .

## docker-run: Run the Docker container
docker-run:
	@echo "==> Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):latest

## help: Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
	@echo ""
	@echo "Pipeline test components: vet -> unit-test -> race-test -> coverage"
