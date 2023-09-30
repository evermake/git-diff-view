#!/usr/bin/env just --justfile

run:
    go run ./cmd/app

build:
    go build -o app ./cmd/app

generate:
    go generate ./...
    oapi-codegen --config oapi.cfg.yaml openapi.yaml

update:
  go get -u
  go mod tidy -v
