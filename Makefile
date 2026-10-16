.PHONY: check fmt vet harness-guard test red-first

check: fmt vet harness-guard test red-first

fmt:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "gofmt: these files need formatting (run gofmt -w .):"; echo "$$out"; exit 1; fi

vet:
	go vet -tags integration ./...

test:
	go test -tags integration ./...

red-first:
	@scripts/red-first.sh

harness-guard:
	@scripts/harness-guard.sh
