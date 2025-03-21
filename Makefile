# --- PREAMBLE
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables

lint-server:
	cd server
	golangci-lint run
.PHONY: lint-server
