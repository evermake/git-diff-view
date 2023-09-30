#!/usr/bin/env just --justfile

generate:
    go generate ./...
    oapi-codegen --config oapi.cfg.yaml openapi.yaml

update:
  go get -u
  go mod tidy -v