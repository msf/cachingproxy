# shamelessly copied from msf/cachingproxy/
APPLICATION := $(shell basename $$PWD)

WORKSPACE := $(realpath ./)
ifeq (,$(wildcard $(WORKSPACE)/bin/.gitkeep))
	WORKSPACE = $(realpath $$PWD)
endif

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
WORKSPACE_BIN := $(WORKSPACE)/bin
WORKSPACE_TMP := $(WORKSPACE)/tmp
WORKSPACE_TMP_BIN := $(WORKSPACE)/tmp/bin
WORKSPACE_TMP_BUILD := $(WORKSPACE)/tmp/build

$(shell mkdir -p $(WORKSPACE_BIN))
$(shell mkdir -p $(WORKSPACE_TMP))
$(shell mkdir -p $(WORKSPACE_TMP_BIN))
$(shell mkdir -p $(WORKSPACE_TMP_BUILD))

export PATH := $(WORKSPACE_TMP_BIN):$(WORKSPACE_BIN):$(PATH)
export GOBIN := $(WORKSPACE_TMP_BIN)

VERSION = $(shell git describe --tags --always --dirty)

# Full IMAGE name, requires the IMAGE_REGISTRY environment variable
IMAGE_REGISTRY = $$IMAGE_REGISTRY
ifndef $$IMAGE_REGISTRY
	IMAGE_REGISTRY = 506714715093.dkr.ecr.us-east-1.amazonaws.com
endif
IMAGE := $(IMAGE_REGISTRY)/$(APPLICATION)/$(APPLICATION):$(VERSION)

# Go compile flags
GOLDFLAGS += -w -extldflags "-static"
GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Name=$(APPLICATION)
GOLDFLAGS += -X main.BuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ%z')
GOMOD = vendor
GOFLAGS = -mod=$(GOMOD) -ldflags "$(GOLDFLAGS)"

# golangci-lint is used as linter
GOLANGCI_VERSION := v1.43.0

# buf.build is used for protobuf and grpc maintenance
PROTO_PATH := $$PWD/proto
BUF_VERSION := 0.56.0

.PHONY: all vendor setup build build-image test-unit test-race lint generate benchmark

all: vendor setup generate build

mod-check:
	@echo "# Checking consistency between go.mod and vendored code..."
	go list -mod=$(GOMOD) ./... >/dev/null
	go mod verify

vendor:
	@echo "# Tidying up go.mod and vendored code..."
	go mod tidy -v
	go mod vendor -v

setup: mod-check 
	@echo "# Workspace $(WORKSPACE)"
ifneq (,$(wildcard ../../tools.go))
	@echo "# Installing tools from global tools.go..."
	@cat ../../tools.go | grep _ | awk -F'"' '{print "go install -v -mod=$(GOMOD) -v "$$2}'|sh -v
endif

ifneq (,$(wildcard ./tools.go))
	@echo "# Installing tools from local tools.go..."
	@cat tools.go | grep _ | awk -F'"' '{print "go install -v -mod=$(GOMOD) -v "$$2}'|sh -v
endif

build: mod-check test-unit
	@echo "# Building target to $(WORKSPACE_TMP_BUILD)/$(APPLICATION)..."
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(WORKSPACE_TMP_BUILD)/$(APPLICATION) ./main.go


image-build:
	@echo "# Building docker image $(IMAGE)..."
	docker build --build-arg LOC=$(notdir $(realpath ../)) \
		--build-arg APPLICATION=$(APPLICATION) -t $(IMAGE) -f Dockerfile ../../


image-push:
	@echo "# Pushing docker image $(IMAGE)..."
	aws --profile unbabel ecr get-login-password --region us-east-1 | \
	docker login --username AWS --password-stdin $(IMAGE_REGISTRY)
	docker push $(IMAGE)


test-unit: mod-check lint
	@echo "# Running unit tests..."
	go test -v -timeout 5m -race -cover ./...

benchmark: mod-check
	@echo "# Running go benchmarks..."
	go test -v -timeout 5m -run=XXX -benchmem -bench ./...

lint: mod-check
	@echo "# Running linters..."
	go vet ./...
ifneq (,$(wildcard ./golangci.yaml))
	@echo "# Running golangci-lint with local config"
	golangci-lint run -c golangci.yaml
else
ifneq (,$(wildcard $(WORKSPACE)/golangci.yaml))
	@echo "# Running golangci-lint with global config"
	golangci-lint run -c $(WORKSPACE)/golangci.yaml
else
	@echo "# Running golangci-lint with default config"
	golangci-lint run
endif
endif

ifneq (,$(wildcard proto/.))
	@echo "# Running buf lint..."
	buf lint
endif


generate: mod-check
	@echo "# Running go generate..."
	go generate ./...
ifneq (,$(wildcard proto/.))
	@echo "# Running buf generate..."
	buf generate
	@echo "# Removing unecessary generated code"
	rm -r proto/gen/go/google
	rm -r proto/gen/openapiv2/google
endif
