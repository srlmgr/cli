package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/srlmgr/cli/log"
)

// NewAPITokenInterceptor returns a client interceptor that sets the Authorization
// header on every outgoing request when token is non-empty.
func NewAPITokenInterceptor(token string) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		//nolint:whitespace // editor/linter issue
		return func(ctx context.Context, req connect.AnyRequest) (
			connect.AnyResponse, error,
		) {
			if token != "" {
				req.Header().Set("api-token", token)
			}
			return next(ctx, req)
		}
	})
}

// NewTraceIDInterceptor returns a client interceptor that examines the response
// for an x-trace-id header and logs it via debug if present.
func NewTraceIDInterceptor(logger *log.Logger) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		//nolint:whitespace // editor/linter issue
		return func(ctx context.Context, req connect.AnyRequest) (
			connect.AnyResponse, error,
		) {
			res, err := next(ctx, req)
			if res != nil {
				if traceID := res.Header().Get("x-trace-id"); traceID != "" {
					if logger != nil {
						logger.Debug("trace id from response",
							log.String("x-trace-id", traceID))
					}
				}
			}
			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				if logger != nil {
					traceID := connectErr.Meta().Get("x-trace-id")
					logger.Error("connect error",
						log.String("error", connectErr.Error()),
						log.Uint32("code", uint32(connectErr.Code())),
						log.String("x-trace-id", traceID),
					)
				}
			}
			return res, err
		}
	})
}
