package client

import (
	"net/http"
	"strings"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/backend/query/v1/queryv1connect"
	"connectrpc.com/connect"

	"github.com/srlmgr/cli/cmd/rpc"
	"github.com/srlmgr/cli/log"
)

// NewQueryServiceClient creates a new authenticated QueryService client.
//
//nolint:whitespace // editor/linter issue
func NewQueryServiceClient(
	apiBaseURL string,
	logger *log.Logger,
) queryv1connect.QueryServiceClient {
	return queryv1connect.NewQueryServiceClient(
		http.DefaultClient,
		normalizeAPIBaseURL(apiBaseURL),
		connect.WithGRPC(),
		connect.WithInterceptors(
			rpc.NewTraceIDInterceptor(logger),
		),
	)
}

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
