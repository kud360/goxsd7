GO ?= go

.PHONY: all build test vet lint check conformance ratchet specs fetch-specs clean

all: check

build:
	$(GO) build ./...

test:
	$(GO) test ./...

# Static checks: go vet + the STYLE.md lint gate (.golangci.yml, which
# also enforces gofmt via its formatters section).
vet:
	$(GO) vet ./...
	golangci-lint run

lint: vet

check: build test vet

# Run the W3C suite against committed expectations; fails on any regression.
conformance:
	$(GO) test ./conformance -run TestConformance -count=1

# Re-baseline expectations UPWARD ONLY (refuses regressions). Arbiter-only.
ratchet:
	GOXSD_RATCHET=1 $(GO) test ./conformance -run TestConformance -count=1

# Regenerate docs/specs/md from the committed HTML.
specs:
	$(GO) run ./tools/spec2md -in docs/specs/html -out docs/specs/md

fetch-specs:
	./scripts/fetch-specs.sh

clean:
	rm -rf bin .agent/logs
