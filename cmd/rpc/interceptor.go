package rpc

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect/v2"

	"github.com/srlmgr/cli/log"
)

// NewAPITokenInterceptor returns a client interceptor that sets the api-token
// header on every outgoing request when token is non-empty.
func NewAPITokenInterceptor(token string) connect.ClientInterceptor {
	return func(next connect.ClientFunc) connect.ClientFunc {
		//nolint:whitespace // editor/linter issue
		return func(ctx context.Context, spec connect.Spec) (
			connect.ClientStream, error,
		) {
			if token != "" {
				if info, ok := connect.CallInfoForClientContext(ctx); ok {
					info.RequestHeader().Set("api-token", token)
				}
			}
			return next(ctx, spec)
		}
	}
}

// NewTraceIDInterceptor returns a client interceptor that examines the response
// for an x-trace-id header and logs it via debug if present. Connect errors are
// logged together with the trace id (taken from response headers or trailers).
func NewTraceIDInterceptor(logger *log.Logger) connect.ClientInterceptor {
	return func(next connect.ClientFunc) connect.ClientFunc {
		//nolint:whitespace // editor/linter issue
		return func(ctx context.Context, spec connect.Spec) (
			connect.ClientStream, error,
		) {
			info, _ := connect.CallInfoForClientContext(ctx)
			stream, err := next(ctx, spec)
			if err != nil {
				logConnectError(logger, info, err)
				return nil, err
			}
			return &traceIDStream{ClientStream: stream, info: info, logger: logger}, nil
		}
	}
}

// traceIDStream inspects response metadata once the first Receive completes.
type traceIDStream struct {
	connect.ClientStream
	info     *connect.CallInfo
	logger   *log.Logger
	reported bool
}

func (s *traceIDStream) Receive(msg any) error {
	err := s.ClientStream.Receive(msg)
	if s.reported || errors.Is(err, io.EOF) {
		return err
	}
	s.reported = true
	if err != nil {
		logConnectError(s.logger, s.info, err)
		return err
	}
	if s.logger != nil && s.info != nil {
		if traceID := s.info.ResponseHeader().Get("x-trace-id"); traceID != "" {
			s.logger.Debug("trace id from response",
				log.String("x-trace-id", traceID))
		}
	}
	return nil
}

func logConnectError(logger *log.Logger, info *connect.CallInfo, err error) {
	var connectErr *connect.Error
	if logger == nil || !errors.As(err, &connectErr) {
		return
	}
	traceID := ""
	if info != nil {
		traceID = info.ResponseHeader().Get("x-trace-id")
		if traceID == "" {
			traceID = info.ResponseTrailer().Get("x-trace-id")
		}
	}
	logger.Error("connect error",
		log.String("error", connectErr.Error()),
		log.Uint32("code", uint32(connectErr.Code())),
		log.String("x-trace-id", traceID),
	)
}
