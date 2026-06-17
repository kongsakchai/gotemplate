.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./... | ./.script/colorize
	@echo "Tests completed."

.PHONY: testcover
testcover:
	@echo "Running tests..."
	@go test -cover ./... | ./.script/colorize
	@echo "Tests completed."

.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./... | ./.script/colorize
	@go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.html"

.PHONY: init
init:
	@chmod +x ./.script/colorize
	@go install github.com/vektra/mockery/v3@v3.7.3
	@go install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: genmock
gen-mock:
	@mockery

.PHONY: gendocs
gendocs:
	@swagger generate spec -o ./docs/swagger.yaml --scan-models --tags=docs

.PHONY: docs
docs:
	@swagger serve -F=swagger ./docs/swagger.yaml
