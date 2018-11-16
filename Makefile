BINARY_NAME=fedrampup

all: test

fmt:
	go fmt

build: fmt
	go build -o $(BINARY_NAME) -v

test: fmt
	go test -v ./...

run: fmt
	go build -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
	go clean
	rm -f $(BINARY_NAME)
