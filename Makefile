.PHONY: test cover coverfunc

test:
	@echo "Running tests..."
	go-acc ./... --ignore internal,docs,server,cmd,repository/repositories/mocks,repository/models,repository/repositories

cover:
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.txt

coverfunc:
	@echo "Generating function coverage report..."
	go tool cover -func=coverage.txt
