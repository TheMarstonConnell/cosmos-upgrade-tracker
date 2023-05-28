install: check-go-version
	go install -mod=readonly ../upgrade-tracker

build: check-go-version
	go build -mod=readonly -o build/upgrade-tracker ../upgrade-tracker

PHONY: install build lint check-go-version

# Add check to make sure we are using the proper Go version before proceeding with anything
check-go-version:
	@if ! go version | grep -q "go1.20"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.20 is required for compiling canined. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.20 and retry. Thank you!"; \
		exit 1; \
	fi


###############################################################################
###                                Linting                                  ###
###############################################################################

format-tools:
	go install mvdan.cc/gofumpt@v0.3.1
	gofumpt -l -w .

lint: format-tools
	golangci-lint run

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/jackalLabs/canine-chain
