.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./... | ./.script/colorize
	@echo "Tests completed."

.PHONY: testcover
testcover:
	@echo "Running tests..."
	@go test -v -cover ./... | ./.script/colorize
	@echo "Tests completed."

.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./... | ./.script/colorize
	@go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.html"

.PHONY: migrate
migrate:
	@echo "Running migrations..."
	@go run ./cmd/migrate/main.go
