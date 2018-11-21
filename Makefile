BINARY_NAME=fedrampup
NAME   := scalesec/fedrampup
TAG    := $(shell git rev-parse HEAD)
IMG    := $(NAME):$(TAG)
LATEST := $(NAME):latest

all: test

fmt:
	go fmt

build: fmt
	go build -o $(BINARY_NAME) -v
	docker build -t $(IMG) .

test: fmt
	go test -v ./...

run: fmt
	go build -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
	go clean
	rm -f $(BINARY_NAME)

release: build
	docker tag $(IMG) $(LATEST)
	docker push $(NAME)
