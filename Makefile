.PHONY: test cover coverfunc

test:
	@echo "Running tests..."
	go test ./... -coverprofile=coverage.out

cover:
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out

coverfunc:
	@echo "Generating function coverage report..."
	go tool cover -func=coverage.out
