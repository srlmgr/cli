package client

import (
	"net/http"
	"strings"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/v2/backend/command/v1/commandv1connect"
	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"

	"github.com/srlmgr/cli/cmd/rpc"
	"github.com/srlmgr/cli/log"
)

// normalizeAPIBaseURL ensures the given URL has a scheme and no trailing whitespace.
func normalizeAPIBaseURL(baseURL string) string {
	url := strings.TrimSpace(baseURL)
	if url == "" {
		return "http://localhost:8080"
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// NewCommandServiceClient creates a ConnectRPC CommandServiceClient
// for the given base URL and optional API token.
//
//nolint:whitespace // editor/linter issue
func NewCommandServiceClient(
	baseURL, token string,
	logger *log.Logger,
) commandv1connect.CommandServiceClient {
	return commandv1connect.NewCommandServiceClient(connect.NewClient(
		connecthttp.NewTransport(
			http.DefaultClient,
			normalizeAPIBaseURL(baseURL),
			connecthttp.WithGRPC(),
			connecthttp.WithReadMaxBytes(0),
		),
		rpc.NewAPITokenInterceptor(token),
		rpc.NewTraceIDInterceptor(logger),
	))
}
