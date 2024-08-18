# Makefile for Go project

# Variables
GOCMD = go
GOLINT = golangci-lint
GOFMT =  gofmt
GOTEST = $(GOCMD) test

test:
	$(GOTEST) -v ./...

fmt:
	$(GOFMT) -s -w .

lint:
	$(GOLINT) run	
