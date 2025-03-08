BINARY := terraform-provider-jsonnet
SOURCES := $(wildcard *.go) $(wildcard private/provider/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

test: $(SOURCES)
	go test -v ./...

clean:
	rm -f $(BINARY)

release:
	goreleaser release --rm-dist

docs:
	go tool tfplugindocs generate --rendered-website-dir docs-new

.PHONY: test release clean docs
