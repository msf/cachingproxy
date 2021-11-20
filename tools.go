//go:build tools
// +build tools

package main

// Import the main binary tools required to work within the monorepo.
// The code is vendored and we will always install from it.
//
// They will be installed to monorepo/tmp/bin and used by all Makefile actions.
//
// Versions are pinned and upgrades must be handled explicitely through
// go get <package>@<version>

import (
	//	_ "github.com/bufbuild/buf/cmd/buf"
	//	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking"
	//	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-lint"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	//	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	//	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	//	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	//	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
