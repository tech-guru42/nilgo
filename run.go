// Copyright (c) 2024 The nilgo authors
// Use of this source code is governed by a MIT license found in the LICENSE file.

package nilgo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nil-go/konf"

	"github.com/nil-go/nilgo/config"
	"github.com/nil-go/nilgo/log"
	"github.com/nil-go/nilgo/run"
)

// Run runs application with the given arguments synchronously.
//
// The runner passed in are running parallel without explicit order,
// which means it should not have temporal dependency between each other.
//
// The running can be interrupted if any runner returns non-nil error, or it receives an OS signal.
// It waits all runners return unless it's forcefully killed by OS.
//
// For now, it only can pass one of following types for args:
//   - config.Option
//   - log.Option
//   - run.Option
//   - func(context.Context) error
func Run(args ...any) error {
	var (
		configOpts []config.Option
		logOpts    []log.Option
		runOpts    []run.Option
		runners    []func(context.Context) error
	)
	for _, arg := range args {
		switch opt := arg.(type) {
		case config.Option:
			configOpts = append(configOpts, opt)
		case log.Option:
			logOpts = append(logOpts, opt)
		case run.Option:
			runOpts = append(runOpts, opt)
		case func(context.Context) error:
			runners = append(runners, opt)
		default:
			return fmt.Errorf("unknown argument type: %T", opt) //nolint:goerr113
		}
	}

	// Initialize the global konf.Config.
	cfg, err := config.New(configOpts...)
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}
	konf.SetDefault(cfg)

	// Initialize the global slog.Logger.
	logger := log.New(logOpts...)
	slog.SetDefault(logger)

	runner := run.New(append(runOpts, run.WithPreRun(cfg.Watch))...)
	if err := runner.Run(context.Background(), runners...); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
