# Background color
GREEN  				:= $(shell tput -Txterm setaf 2)
YELLOW 				:= $(shell tput -Txterm setaf 3)
BLUE 				:= $(shell tput -Txterm setaf 4)
MAGENTA             := $(shell tput -Txterm setaf 5)
WHITE  				:= $(shell tput -Txterm setaf 7)
RESET  				:= $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM := 20


## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET} ${MAGENTA}[variable=value]${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR  := $(dir $(MKFILE_PATH))

.PHONY: build
## Build command line interface
build:
	echo "Building"

.PHONY: format
## Format *.go by go format
format:
	go fmt ./...

.PHONY: sformat
## Strictly format *.go by gofumpt
sformat:
	gofumpt -w -s ./

.PHONY: lint
## Lint *.go via golangci-lint
lint:
	golangci-lint run -v

.PHONY: test
## test go files
test:
	go test -failfast -race -v -coverprofile=coverage.out ./...

.PHONY: generate
## Generate generated go files via go genreate tool
generate:
	go generate ./...

.PHONY: bench
## benchmark everything
bench:
	go test -benchmem -run="^$$" -bench ^Benchmark ./...

.PHONY: fuzz
## fuzzying everything
fuzzy:
	echo "fuzzying"

.PHONY: clean
## clean unused files
clean:
	rm *-fuzz.zip *.log *.out *.test bin/*
