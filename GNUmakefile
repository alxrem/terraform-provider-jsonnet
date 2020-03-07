BINARY := terraform-provider-jsonnet
SOURCES := $(wildcard *.go) $(wildcard jsonnet/*.go)

default: $(BINARY)

$(BINARY):  $(SOURCES)
	go build -o $(BINARY)
