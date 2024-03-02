.SILENT :

export GO111MODULE=on

# App name
APPNAME=readflow

# Go configuration
GOOS?=$(shell go env GOHOSTOS)
GOARCH?=$(shell go env GOHOSTARCH)

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

# Extract version infos
PKG_VERSION:=github.com/ncarlier/$(APPNAME)/internal/version
VERSION:=`git describe --always --dirty`
GIT_COMMIT:=`git rev-list -1 HEAD --abbrev-commit`
BUILT:=`date`
define LDFLAGS
-X '$(PKG_VERSION).Version=$(VERSION)' \
-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
-X '$(PKG_VERSION).Built=$(BUILT)'
endef

all: build

# Include common Make tasks
root_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
-include $(root_dir)/.env
makefiles:=$(root_dir)/makefiles
include $(makefiles)/help.Makefile
include $(makefiles)/docker/compose.Makefile

# Some variables
db_service?=readflow-db-1

## Clean built files
clean:
	echo ">>> Removing generated files..."
	-rm -rf release autogen
.PHONY: clean

## Run code generation
autogen:
	echo ">>> Generating code ..."
	-mkdir -p autogen/db/postgres
	go generate

## Build executable
build: autogen
	-mkdir -p release
	echo ">>> Building $(EXECUTABLE) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o release/$(EXECUTABLE)
.PHONY: build

release/$(EXECUTABLE): build

# Check code style
check-style:
	echo ">>> Checking code style..."
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
.PHONY: check-style

# Check code criticity
check-criticity:
	echo ">>> Checking code criticity..."
	go run github.com/go-critic/go-critic/cmd/gocritic@latest check -enableAll ./...
.PHONY: check-criticity

# Check code security
check-security:
	echo ">>> Checking code security..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest -quiet ./...
.PHONY: check-security

## Code quality checks
checks: check-style check-criticity
.PHONY: checks

## Run tests
test:
	echo ">>> Running tests..."
	go test ./...
.PHONY: test

## Run test coverage
test-cov:
	echo ">>> Running test coverage..."
	go test -coverprofile=release/coverage.out ./...
	go tool cover -html=release/coverage.out -o release/coverage.html
.PHONY: test-cov

## Install executable
install: release/$(EXECUTABLE)
	echo ">>> Installing $(EXECUTABLE) to ${HOME}/.local/bin/$(EXECUTABLE) ..."
	cp release/$(EXECUTABLE) ${HOME}/.local/bin/$(EXECUTABLE)
.PHONY: install

## Create Docker image
image:
	echo ">>> Building Docker image..."
	docker build --rm -t ncarlier/$(APPNAME) .
.PHONY: image

# Generate changelog
CHANGELOG.md:
	standard-changelog --first-release

## Generate documentation website
docs:
	echo ">>> Building documentation website..."
	hugo -s docs -d ../release/docs --cleanDestinationDir
.PHONY: docs

## Generate landing pages
landing:
	echo ">>> Building landing pages..."
	cd landing && npm install --silent && npm run build
.PHONY: landing

## Generate Web UI
ui:
	echo ">>> Building Web UI..."
	cd ui && npm install --silent --legacy-peer-deps && REACT_APP_VERSION=${VERSION} npm run build
.PHONY: ui

## Build bookmarklet
bookmarklet:
	echo ">>> Building Bookmarklet..."
	cd bookmarklet && npm install --silent && npm run clean && npm run build
	cp dist/bookmarklet.html ../ui/public/
	cp dist/bookmarklet.*.js ../ui/public/bookmarklet.js
.PHONY: bookmarklet

## Create archive
archive: release/$(EXECUTABLE) CHANGELOG.md
	echo ">>> Creating release/$(ARCHIVE) archive..."
	tar czf release/$(ARCHIVE) README.md LICENSE -C release/ $(EXECUTABLE)
	rm release/$(EXECUTABLE)
.PHONY: archive

## Create distribution binaries
distribution:
	GOARCH=amd64 make build archive
	GOARCH=arm64 make build archive
	GOARCH=arm make build archive
	GOOS=darwin make build archive
	GOOS=windows make build archive
.PHONY: distribution

## Start development server (aka: a test database instance)
dev-server:
	docker compose -f docker-compose.dev.yml down
	docker compose -f docker-compose.dev.yml up
.PHONY: dev-server

## Deploy containers to Docker host
deploy: compose-build compose-up
.PHONY: deploy

## Un-deploy containers from Docker host
undeploy: compose-down
.PHONY: undeploy

## Backup database
backup:
	archive=backup/$(db_service)-`date -I`.dump
	echo "Backuping PosgreSQL database ($(db_service) ==> $$archive)..."
	mkdir -p backup
	docker exec -u postgres $(db_service) pg_dumpall > $$archive
	gzip -f $$archive
	echo "done."
.ONESHELL:
.PHONY: backup

## Restore database
restore:
	echo "Restoring $(archive) database dump to $(db_service) ..."
	@while [ -z "$$CONTINUE" ]; do \
		read -r -p "Are you sure? [y/N]: " CONTINUE; \
	done ; \
	[ $$CONTINUE = "y" ] && docker exec -i -u postgres $(db_service) psql -U postgres -d postgres < $(archive)
.PHONY: restore

## Open database client
db-client:
	docker exec -it -u postgres $(db_service) psql -U postgres
.PHONY: db-client

var/block-list.txt:
	echo ">>> Downloading blocklist file..."
	mkdir -p var
	wget -O var/block-list.txt https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt
