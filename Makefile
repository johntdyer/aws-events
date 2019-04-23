SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=aws-event-monitor

VERSION=0.0.2

LDFLAGS=-ldflags "-X main.Build=`git rev-parse HEAD` -a -installsuffix cgo"

.DEFAULT_GOAL: $(BINARY)

packages/$(BINARY): $(SOURCES)
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./packages/aws-events_v${VERSION}

.PHONY: clean
clean:
		@rm -vf $(BINARY)_*_386 $(BINARY)_*_amd64 $(BINARY)_*_arm $(BINARY)_*.exe
