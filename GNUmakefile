BINARY := terraform-provider-jsonnet
SOURCES := $(wildcard *.go) $(wildcard jsonnet/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

test: $(SOURCES)
	go test ./...

release:
	goreleaser release --rm-dist

.PHONY: test release
