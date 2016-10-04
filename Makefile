ifndef GOPATH
	$(warning You need to set up a GOPATH.)
endif

PROJECT := github.com/juju/bundleservice
PROJECT_DIR := $(shell go list -e -f '{{.Dir}}' $(PROJECT))

GIT_COMMIT := $(shell git rev-parse --verify HEAD)
BUNDLECHANGES_COMMIT := $(shell grep bundlechanges dependencies.tsv | cut -f 3)

help:
	@echo "Available targets:"
	@echo "  deps - fetch all dependencies"
	@echo "  build - build the project"
	@echo "  check - run tests"
	@echo "  install - install the library in your GOPATH"
	@echo "  clean - clean the project"

run: params/init.go
	go run server.go

# Start of GOPATH-dependent targets. Some targets only make sense -
# and will only work - when this tree is found on the GOPATH.
ifeq ($(CURDIR),$(PROJECT_DIR))

deps: params/init.go
	go get -v -t $(PROJECT)/...

build: params/init.go
	go build $(PROJECT)/...

check: params/init.go
	go test $(PROJECT)/...

install:
	go install $(INSTALL_FLAGS) -v $(PROJECT)/...

clean:
	go clean $(PROJECT)/...

# Generate version information
params/init.go: params/init.go.tmpl
	gofmt -r "unknownVersion -> VersionInfo{GitCommit: \"${GIT_COMMIT}\", BundlechangesCommit: \"${BUNDLECHANGES_COMMIT}\",}" $< > $@

else

run:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

deps:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

build:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

check:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

install:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

clean:
	$(error Cannot $@; $(CURDIR) is not on GOPATH)

endif
# End of GOPATH-dependent targets.

.PHONY: help deps build check install clean
