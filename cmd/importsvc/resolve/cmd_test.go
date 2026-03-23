//nolint:lll,whitespace // test code
package resolve

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	"connectrpc.com/connect"
)

// ---- mock implementations ----

type mockImportClient struct {
	resolveMappings func(context.Context, *connect.Request[importv1.ResolveMappingsRequest]) (*connect.Response[importv1.ResolveMappingsResponse], error)
}

func (m *mockImportClient) ResolveMappings(
	ctx context.Context,
	req *connect.Request[importv1.ResolveMappingsRequest],
) (*connect.Response[importv1.ResolveMappingsResponse], error) {
	if m.resolveMappings != nil {
		return m.resolveMappings(ctx, req)
	}

	resp := &importv1.ResolveMappingsResponse{}
	resp.SetUnresolvedMappings(3)

	return connect.NewResponse(resp), nil
}

// ---- tests ----

func TestResolveCommand_Success(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &resolveCommand{
		raceID:    7,
		out:       &buf,
		importSvc: &mockImportClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "unresolved_mappings=3") {
		t.Errorf("expected output to contain unresolved_mappings=3, got: %s", out)
	}
}

func TestResolveCommand_PassesRaceID(t *testing.T) {
	t.Parallel()

	var capturedRaceID uint32

	runner := &resolveCommand{
		raceID: 42,
		out:    &bytes.Buffer{},
		importSvc: &mockImportClient{
			resolveMappings: func(_ context.Context, req *connect.Request[importv1.ResolveMappingsRequest]) (*connect.Response[importv1.ResolveMappingsResponse], error) {
				capturedRaceID = req.Msg.GetRaceId()
				resp := &importv1.ResolveMappingsResponse{}
				resp.SetUnresolvedMappings(0)
				return connect.NewResponse(resp), nil
			},
		},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedRaceID != 42 {
		t.Errorf("expected race_id=42, got %d", capturedRaceID)
	}
}

func TestResolveCommand_ServiceError(t *testing.T) {
	t.Parallel()

	runner := &resolveCommand{
		raceID: 1,
		out:    &bytes.Buffer{},
		importSvc: &mockImportClient{
			resolveMappings: func(_ context.Context, _ *connect.Request[importv1.ResolveMappingsRequest]) (*connect.Response[importv1.ResolveMappingsResponse], error) {
				return nil, errors.New("service error")
			},
		},
	}

	err := runner.run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "resolve mappings") {
		t.Errorf("expected error to mention 'resolve mappings', got: %v", err)
	}
}
