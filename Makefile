SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=aws-event-monitor

VERSION=0.0.1

LDFLAGS=-ldflags "-X main.Build=`git rev-parse HEAD`"

.DEFAULT_GOAL: $(BINARY)

packages/$(BINARY): $(SOURCES)
		gox \
		-osarch="linux/amd64"    \
		-osarch="darwin/amd64" \
		${LDFLAGS}               \
		-output packages/{{.Dir}}_{{.OS}}_v${VERSION}_{{.Arch}}

.PHONY: clean
clean:
		@rm -vf $(BINARY)_*_386 $(BINARY)_*_amd64 $(BINARY)_*_arm $(BINARY)_*.exe
