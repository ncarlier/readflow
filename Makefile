.SILENT :

export GO111MODULE=on

# App name
APPNAME=readflow

# Go configuration
GOOS?=linux
GOARCH?=amd64

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

# Extract version infos
PKG_VERSION:=github.com/ncarlier/$(APPNAME)/pkg/version
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

## Run tests
test:
	echo ">>> Running tests..."
	go test ./...
.PHONY: test

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

## Generate changelog
changelog:
	standard-changelog --first-release
.PHONY: changelog

## Generate HTML website
html:
	echo ">>> Building static website..."
	hugo -s landing -d ../release/html --cleanDestinationDir
	hugo -s docs -d ../release/html/docs --cleanDestinationDir
.PHONY: html

## Create archive
archive: release/$(EXECUTABLE)
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
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml up
.PHONY: dev-server

## Start test server (aka: full stack service with mocks)
test-server:
	docker-compose -f docker-compose.test.yml down
	docker-compose -f docker-compose.test.yml up
.PHONY: test-server

## Deploy containers to Docker host
deploy: compose-build compose-up
.PHONY: deploy

## Un-deploy containers from Docker host
undeploy: compose-down
.PHONY: undeploy

## Backup database
backup:
	archive=backup/db-`date -I`.dump
	echo "Backuping PosgreSQL database ($$archive)..."
	mkdir -p backup
	docker exec -u postgres nxreaderapi_database_1 pg_dump -Fc keycloak > $$archive
	gzip $$archive
	echo "done."
.ONESHELL:
.PHONY: backup

## Restore database
restore:
	echo "Restoring $(archive) database dump ..."
	@while [ -z "$$CONTINUE" ]; do \
		read -r -p "Are you sure? [y/N]: " CONTINUE; \
	done ; \
	[ $$CONTINUE = "y" ] || [ $$CONTINUE = "Y" ] || (echo "Exiting."; exit 1;)
	docker exec -i -u postgres nxreaderapi_database_1 pg_restore -C -d postgres < $(archive)
.PHONY: restore
