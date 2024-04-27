// Copyright (c) 2024 The nilgo authors
// Use of this source code is governed by a MIT license found in the LICENSE file.

package grpc

import (
	_ "unsafe" // For go:linkname

	"google.golang.org/grpc"
)

//go:linkname addGlobalDialOptions google.golang.org/grpc/internal.AddGlobalDialOptions
var addGlobalDialOptions any //nolint:gochecknoglobals // func(opt ...DialOption)

// WithGlobalDialOption adds global dial options for all gRPC clients.
//
// CAUTION: This function may break in new version of `google.golang.org/grpc`
// since it is using internal package from grpc.
func WithGlobalDialOption(opts ...grpc.DialOption) {
	addGlobalDialOptions.(func(opt ...grpc.DialOption))(opts...) //nolint:forcetypeassert
}
