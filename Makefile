# --- PREAMBLE
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables

format-client:
	cd client
	pnpm format
.PHONY: format-client

lint: lint-client lint-server
.PHONY: lint

lint-client:
	cd client
	pnpm lint
.PHONY: lint-client

lint-server:
	cd server
	golangci-lint run
.PHONY: lint-server
