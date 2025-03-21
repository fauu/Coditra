# --- PREAMBLE
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables

build: build-client build-server
.PHONY: build

build-client:
	cd client
	pnpm install
	pnpm build
.PHONY: build-server

build-server:
	./scripts/server-prebuild.sh
	cd server
	go build -o target/coditra cmd/coditra/main.go
.PHONY: build-server

dev-client:
	cd client
	pnpm dev
.PHONY: dev-client

dev-server:
	cd server
	air
.PHONY: dev-server

format-client:
	cd client
	pnpm format
.PHONY: format-client

lint: lint-client lint-server
.PHONY: lint

lint-client:
	cd client
	pnpm lint
	pnpm svelte-check
.PHONY: lint-client

lint-server:
	cd server
	golangci-lint run
.PHONY: lint-server
