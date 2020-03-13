BINARY := terraform-provider-jsonnet
SOURCES := $(wildcard *.go) $(wildcard jsonnet/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

.PHONY: test
test: $(SOURCES)
	go test ./...
