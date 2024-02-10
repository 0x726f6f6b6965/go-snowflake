PROJECTNAME := $(shell basename "$(PWD)")

## help: Print usage information
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECTNAME)"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo

## test-go: Test go file and show the coverage
.PHONY: test-go
test-go:
	@go test --coverprofile=coverage.out ./... 
	@go tool cover -html=coverage.out  