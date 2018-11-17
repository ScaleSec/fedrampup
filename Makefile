SHELL := /bin/bash

REV := $(shell git rev-parse HEAD)
CHANGES := $(shell test -n "$$(git status --porcelain)" && echo '+CHANGES' || true)

TARGET := fedrampup
VERSION := $(shell cat VERSION)

OS := darwin freebsd linux openbsd
ARCH := 386 amd64
LDFLAGS := -X github.com/ScaleSec/$(TARGET)/local.Revision=$(REV)$(CHANGES)

GPG_SIGNING_KEY :=

.PHONY: \
	help \
	default \
	clean \
	clean-artifacts \
	clean-releases \
	clean-vendor \
	tools \
	deps \
	test \
	coverage \
	vet \
	lint \
	imports \
	fmt \
	env \
	build \
	build-all \
	doc \
	release \
	package-release \
	sign-release \
	check \
	vendor \
	version

all: imports fmt lint vet build

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binaries, artifacts and releases.'
	@echo '    clean-artifacts    Remove build artifacts only.'
	@echo '    clean-releases     Remove releases only.'
	@echo '    clean-vendor       Remove content of the vendor directory.'
	@echo '    tools              Install tools needed by the project.'
	@echo '    deps               Download and install build time dependencies.'
	@echo '    test               Run unit tests.'
	@echo '    coverage           Report code tests coverage.'
	@echo '    vet                Run go vet.'
	@echo '    lint               Run golint.'
	@echo '    imports            Run goimports.'
	@echo '    fmt                Run go fmt.'
	@echo '    env                Display Go environment.'
	@echo '    build              Build project for current platform.'
	@echo '    build-all          Build project for all supported platforms.'
	@echo '    doc                Start Go documentation server on port 8080.'
	@echo '    release            Package and sing project for release.'
	@echo '    package-release    Package release and compress artifacts.'
	@echo '    sign-release       Sign release and generate checksums.'
	@echo '    check              Verify compiled binary.'
	@echo '    vendor             Update and save project build time dependencies.'
	@echo '    version            Display Go version.'
	@echo ''
	@echo 'Targets run by default are: imports, fmt, lint, vet, and build.'
	@echo ''

print-%:
	@echo $* = $($*)

clean: clean-artifacts clean-releases
	go clean -i ./...
	rm -vf \
	  $(CURDIR)/coverage.* \

clean-artifacts:
	rm -Rf artifacts/*

clean-releases:
	rm -Rf releases/*

clean-vendor:
	find $(CURDIR)/vendor -type d -print0 2>/dev/null | xargs -0 rm -Rf

clean-all: clean clean-artifacts clean-vendor

tools:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint
	go get github.com/axw/gocov/gocov
	go get github.com/matm/gocov-html
	go get github.com/tools/godep
	go get github.com/mitchellh/gox

deps:
	godep restore

test: deps
	go test -v ./...

coverage: deps
	gocov test ./... > $(CURDIR)/coverage.out 2>/dev/null
	gocov report $(CURDIR)/coverage.out
	if test -z "$$CI"; then \
	  gocov-html $(CURDIR)/coverage.out > $(CURDIR)/coverage.html; \
	  if which open &>/dev/null; then \
	    open $(CURDIR)/coverage.html; \
	  fi; \
	fi

vet:
	go vet -v ./...

lint:
	golint ./...

imports:
	goimports -l -w .

fmt:
	go fmt ./...

env:
	@go env

build: deps
	go build -v \
	   -ldflags "$(LDFLAGS)" \
	   -o "$(TARGET)" .

build-all: deps
	mkdir -v -p $(CURDIR)/artifacts/$(VERSION)
	gox -verbose \
	    -os "$(OS)" -arch "$(ARCH)" \
	    -ldflags "$(LDFLAGS)" \
	    -output "$(CURDIR)/artifacts/$(VERSION)/{{.OS}}_{{.Arch}}/$(TARGET)" .
	cp -v -f \
	   $(CURDIR)/artifacts/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)/$(TARGET) .

doc:
	godoc -http=:8080 -index

release: package-release sign-release

package-release:
	@test -x $(CURDIR)/artifacts/$(VERSION) || exit 1
	mkdir -v -p $(CURDIR)/releases/$(VERSION)
	for release in $$(find $(CURDIR)/artifacts/$(VERSION) -mindepth 1 -maxdepth 1 -type d 2>/dev/null); do \
	  platform=$$(basename $$release); \
	  pushd $$release &>/dev/null; \
	  zip $(CURDIR)/releases/$(VERSION)/$(TARGET)_$${platform}.zip $(TARGET); \
	  popd &>/dev/null; \
	done

sign-release:
	@test -x $(CURDIR)/releases/$(VERSION) || exit 1
	pushd $(CURDIR)/releases/$(VERSION) &>/dev/null; \
	shasum -a 256 -b $(TARGET)_* > SHA256SUMS; \
	if test -n "$(GPG_SIGNING_KEY)"; then \
	  gpg --default-key $(GPG_SIGNING_KEY) -a \
	      -o SHA256SUMS.sign -b SHA256SUMS; \
	fi; \
	popd &>/dev/null

check:
	@test -x $(CURDIR)/$(TARGET) || exit 1
	if $(CURDIR)/$(TARGET) --version | grep -qF '$(VERSION)'; then \
	  echo "$(CURDIR)/$(TARGET): OK"; \
	else \
	  exit 1; \
	fi

vendor: deps
	godep save

version:
	@go version

