//go:build tools

package cmd

import (
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/golang/mock/gomock"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/xo/xo"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "gotest.tools/gotestsum"
)
