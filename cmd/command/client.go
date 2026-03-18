package command

import (
	"net/http"
	"strings"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/backend/command/v1/commandv1connect"
	"connectrpc.com/connect"

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
	return commandv1connect.NewCommandServiceClient(
		http.DefaultClient,
		normalizeAPIBaseURL(baseURL),
		connect.WithGRPC(),
		connect.WithInterceptors(
			rpc.NewAPITokenInterceptor(token),
			rpc.NewTraceIDInterceptor(logger),
		),
	)
}
