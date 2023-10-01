#!/usr/bin/env just --justfile

app := "git-diff-server"

build:
    go build -o {{ app }} ./cmd/app

install path: build
    cp {{ app }} {{ path }}
    rm {{ app }}

clean:
    rm {{ app }}

generate:
    go generate ./...
    oapi-codegen --config oapi.cfg.yaml openapi.yaml

update:
  go get -u
  go mod tidy -v
