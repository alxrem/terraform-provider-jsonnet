BINARY := terraform-provider-jsonnet
SOURCES := $(wildcard *.go) $(wildcard jsonnet/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

test: $(SOURCES)
	go test ./...

clean:
	rm -f $(BINARY)

release:
	goreleaser release --rm-dist

.PHONY: test release
