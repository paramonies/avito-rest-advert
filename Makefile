# Go related commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...

.PHONY: build
build: 
	${GOBUILD} -v ./cmd/apiserver

.PHONY: test
test:
	$(GOTEST)

# Generates a coverage report
.PHONY: cover
cover:
	${GOCMD} test -coverprofile=coverage.out ./... && ${GOCMD} tool cover -html=coverage.out

.PHONY: migrate
migrate:
	migrate -path ./migrations -database postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable up

.PHONY: swag
swag:
	swag init -g ./cmd/apiserver/main.go 

.PHONY: clean
clean:
	@rm -f coverage.out

.DEFAULT_GOAL := build