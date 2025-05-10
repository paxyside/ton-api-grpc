.PHONY: all, prepare, lint-fix, lint, pack, tests, local_run, clear, docker_run, proto, setup

SHELL := /bin/bash
.SHELLFLAGS := -e -c

main_path = ./cmd/server/main.go
app_name = ./server
proto_src_dir = proto
proto_out_dir = .
dc = docker compose

all: clear lint pack local_run
prepare: lint-fix tests proto

lint:
	go tool golangci-lint run

lint-fix:
	go tool golangci-lint run --fix

pack:
	go build -o $(app_name) $(main_path)

tests:
	go test -v ./...

local_run:
	$(app_name)

clear:
	rm -f $(app_name)

docker_run:
	$(dc) up --build -d

proto:
	protoc --proto_path=$(proto_src_dir) --go_out=$(proto_out_dir) --go-grpc_out=$(proto_out_dir)  $(proto_src_dir)/tonnode/ton_node.proto

setup:
	cp config.example.yaml config.yaml
	docker network create ton-node-network
