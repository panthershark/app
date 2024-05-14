#! /bin/bash

go mod tidy
go mod vendor
podman compose -f test-compose.yaml down --remove-orphans
podman image rm backend_test_backend:latest
podman-compose -f test-compose.yaml up --build --abort-on-container-exit --exit-code-from test_backend

