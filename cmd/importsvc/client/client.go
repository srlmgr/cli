package client

import (
	"net/http"
	"strings"

	"buf.build/gen/go/srlmgr/api/connectrpc/go/v2/backend/import/v1/importv1connect"
	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"

	"github.com/srlmgr/cli/cmd/rpc"
	"github.com/srlmgr/cli/log"
)

// NewImportServiceClient creates a new authenticated ImportService client.
//
//nolint:whitespace // editor/linter issue
func NewImportServiceClient(
	apiBaseURL, token string,
	logger *log.Logger,
) importv1connect.ImportServiceClient {
	return importv1connect.NewImportServiceClient(connect.NewClient(
		connecthttp.NewTransport(
			http.DefaultClient,
			normalizeAPIBaseURL(apiBaseURL),
			connecthttp.WithGRPC(),
			connecthttp.WithReadMaxBytes(0),
		),
		rpc.NewAPITokenInterceptor(token),
		rpc.NewTraceIDInterceptor(logger),
	))
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
