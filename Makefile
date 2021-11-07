# CGO_ENABLED=0 is used becuase we don't want to rely on c libs and opentelemetry
# also builds their binaries this way.
# ref: https://github.com/open-telemetry/opentelemetry-collector/blob/4c503ddc/Makefile#L254-L256
CGO_ENABLED ?= 0
GO ?= go
VERSION ?= "$(shell git describe --tags --abbrev=10)"
BINARY_NAME ?= collector
BINARY_PATH = bin/$(BINARY_NAME)
GENERATE_SOURCES_PATH ?= cmd/
CORE_VERSION ?= 0.37.0

# ALL_MODULES includes ./* dirs (excludes . dir and example with go code)
ALL_MODULES := $(shell find ./pkg -type f -name "go.mod" -exec dirname {} \; | sort | egrep  '^./' )

.PHONY: _generate
_generate:
	CGO_ENABLED=$(CGO_ENABLED) opentelemetry-collector-builder \
		--go $(GO) \
		--version $(VERSION) \
		--skip-compilation=true \
		--output-path ./$(GENERATE_SOURCES_PATH) \
		--config builder_config.yaml

.PHONY: _build
_build:
	(cd $(GENERATE_SOURCES_PATH) && \
		CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="-s -w" \
		-trimpath -o ../$(BINARY_PATH) . \
	) && chmod +x ./$(BINARY_PATH)

.PHONY: build
build: _generate _build

.PHONY: install-builder
install-builder:
	go install github.com/open-telemetry/opentelemetry-collector-builder@v$(CORE_VERSION)

.PHONY: gomod-download-all
gomod-download-all:
	@$(MAKE) for-all CMD="make mod-download-all"

.PHONY: for-all
for-all:
	@echo "running $${CMD} in all modules..."
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
		echo "running $${CMD} in $${dir}" && \
		$${CMD} ); \
	done
