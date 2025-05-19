PHONY: gin
gin:
	@echo "Starting Example Gin server..."
	@go run ./example/gin

PHONY: echo
echo:
	@echo "Starting Example Echo server..."
	@go run ./example/echo


PHONY: test
test:
	@echo "Running tests..."
	@go test -v -tags=test ./... 
	@echo "Tests completed."

PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -tags=test -coverprofile=coverage.out ./... 
	@go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.html"
