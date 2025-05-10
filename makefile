PHONY: gin
gin:
	@echo "Starting Example Gin server..."
	@go run ./example/gin

PHONY: echo
echo:
	@echo "Starting Example Echo server..."
	@go run ./example/echo
