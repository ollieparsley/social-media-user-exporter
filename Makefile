# MAKEFILE
#
# @author      Ollie Parsley <ollie@ollieparsley.com>
# @link        https://github.com/ollieparsley/social-media-user-exporter
# ------------------------------------------------------------------------------

# Ensure everyone is using bash. Note that Ubuntu now uses dash which doesn't support PIPESTATUS.
SHELL=/bin/bash

# Project version
MINOR_VERSION = 0
ifdef GITHUB_RUN_NUMBER
  MINOR_VERSION = $(GITHUB_RUN_NUMBER)
endif
VERSION=$(shell cat VERSION).$(MINOR_VERSION)

# Name of package
PKGNAME=social-media-user-exporter

# Destdir
DESTDIR?=target/build/root/

# Binary path (where the executable files will be installed)
BINPATH=usr/bin/

# Docker
DOCKER_USERNAME?=ollieparsley
DOCKER_PASSWORD?=not_this_password

# Environment variable exports
export GO111MODULE=on

# Init
init:
	go mod init

# Get the dependencies
deps:
	go get ./...
	go get golang.org/x/lint/golint

# Run the unit tests
test:
	go test -v ./...

# Check for syntax errors
vet:
	go vet ./...

# Go fmt
fmt:
	go fmt ./...

# Check for style errors
lint:
	golint -set_exit_status ./...

# Alias to run targets: fmtcheck test vet lint
qa: fmt test vet lint

# Compile the application
build: deps
	@mkdir -p $(DESTDIR)
	@mkdir -p $(DESTDIR)$(BINPATH)
	go build -o $(DESTDIR)$(BINPATH)$(PKGNAME) ./main.go

# Compile the application
run:
	SMUE_INTERVAL_SECONDS="300" \
	SMUE_TWITTER_SCREEN_NAMES="ollieparsley,olliedude2k,meltwatereng" \
	SMUE_TWITTER_ACCESS_TOKEN="$(shell cat resources/env/twitter/access_token)" \
	SMUE_TWITTER_ACCESS_TOKEN_SECRET="$(shell cat resources/env/twitter/access_token_secret)" \
	SMUE_TWITTER_CLIENT_ID="$(shell cat resources/env/twitter/client_id)" \
	SMUE_TWITTER_CLIENT_SECRET="$(shell cat resources/env/twitter/client_secret)" \
	SMUE_YOUTUBE_CHANNEL_IDS="UC7R4lEiVaathpWrwXArnKQg" \
	SMUE_YOUTUBE_CLIENT_ID="$(shell cat resources/env/youtube/client_id)" \
	SMUE_YOUTUBE_CLIENT_SECRET="$(shell cat resources/env/youtube/client_secret)" \
	SMUE_YOUTUBE_ACCESS_TOKEN="$(shell cat resources/env/youtube/access_token)" \
	SMUE_YOUTUBE_REFRESH_TOKEN="$(shell cat resources/env/youtube/refresh_token)" \
	go run ./main.go

# Docker tags
docker-tags:
	echo "latest,$(VERSION)" > .tags