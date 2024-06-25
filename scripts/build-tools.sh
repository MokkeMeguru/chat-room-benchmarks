#!/usr/bin/env sh

set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "${ROOT_DIR}/tools/cmd"

go mod download
go mod tidy

export CGO_ENABLED=0

go build -o "${ROOT_DIR}/bin/protoc-gen-connect-go" "connectrpc.com/connect/cmd/protoc-gen-connect-go"
go build -o "${ROOT_DIR}/bin/buf" "github.com/bufbuild/buf/cmd/buf"
go build -o "${ROOT_DIR}/bin/oapi-codegen" "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
go build -o "${ROOT_DIR}/bin/grpcurl" "github.com/fullstorydev/grpcurl/cmd/grpcurl"
go build -o "${ROOT_DIR}/bin/mockgen" "github.com/golang/mock/mockgen"
go build -o "${ROOT_DIR}/bin/golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint"
go build -o "${ROOT_DIR}/bin/xo" "github.com/xo/xo"
go build -o "${ROOT_DIR}/bin/goimports" "golang.org/x/tools/cmd/goimports"
go build -o "${ROOT_DIR}/bin/protoc-gen-go" "github.com/golang/protobuf/protoc-gen-go"
go build -o "${ROOT_DIR}/bin/gotestsum" "gotest.tools/gotestsum"

# FIXME bin/protoc-gen-go が機能しないので、go install でインストールする
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/bojand/ghz/cmd/ghz@latest
