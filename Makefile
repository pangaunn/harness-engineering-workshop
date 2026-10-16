.PHONY: check fmt vet test

check: fmt vet test

fmt:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "gofmt: these files need formatting (run gofmt -w .):"; echo "$$out"; exit 1; fi

vet:
	go vet -tags integration ./...

test:
	go test -tags integration ./...
