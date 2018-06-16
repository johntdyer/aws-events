SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=aws-event-monitor

VERSION=0.0.1

LDFLAGS=-ldflags "-X main.Build=`git rev-parse HEAD`"

.DEFAULT_GOAL: $(BINARY)

packages/$(BINARY): $(SOURCES)
		gox \
		-osarch="!darwin/386"    \
		-osarch="!linux/386"     \ 
		-osarch="!openbsd/386"   \
		-osarch="!openbsd/amd64" \
		-osarch="!freebsd/amd64" \
		-osarch="!freebsd/arm"   \
		-osarch="!freebsd/386"   \
		-osarch="!windows/386"   \
		-osarch="!windows/amd64" \
		-osarch="!linux/arm"     \
		-osarch="!netbsd/arm"    \
		-osarch="!netbsd/386"    \
		-osarch="!netbsd/amd64"  \
		${LDFLAGS}               \
		-output packages/{{.Dir}}_{{.OS}}_v${VERSION}_{{.Arch}}

.PHONY: clean
clean:
		@rm -vf $(BINARY)_*_386 $(BINARY)_*_amd64 $(BINARY)_*_arm $(BINARY)_*.exe
