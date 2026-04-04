//nolint:lll,whitespace // test code
package preview

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
)

// ---- mock implementations ----

type mockQueryClient struct {
	getRace func(context.Context, *connect.Request[queryv1.GetRaceRequest]) (*connect.Response[queryv1.GetRaceResponse], error)
}

func (m *mockQueryClient) GetRace(
	ctx context.Context,
	req *connect.Request[queryv1.GetRaceRequest],
) (*connect.Response[queryv1.GetRaceResponse], error) {
	if m.getRace != nil {
		return m.getRace(ctx, req)
	}

	race := &commonv1.Race{}
	race.SetId(req.Msg.GetId())
	race.SetEventId(10)

	resp := &queryv1.GetRaceResponse{}
	resp.SetRace(race)

	return connect.NewResponse(resp), nil
}

type mockImportClient struct {
	getPreprocessPreview func(context.Context, *connect.Request[importv1.GetPreprocessPreviewRequest]) (*connect.Response[importv1.GetPreprocessPreviewResponse], error)
}

func (m *mockImportClient) GetPreprocessPreview(
	ctx context.Context,
	req *connect.Request[importv1.GetPreprocessPreviewRequest],
) (*connect.Response[importv1.GetPreprocessPreviewResponse], error) {
	if m.getPreprocessPreview != nil {
		return m.getPreprocessPreview(ctx, req)
	}

	return connect.NewResponse(&importv1.GetPreprocessPreviewResponse{}), nil
}

// ---- tests ----

func TestPreviewCommand_Success(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &previewCommand{
		raceGridID: 1,
		out:        &buf,
		qrySvc:     &mockQueryClient{},
		importSvc:  &mockImportClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "rows=0") {
		t.Errorf("expected output to contain rows=0, got: %s", out)
	}
	if !strings.Contains(out, "unresolved_mappings=0") {
		t.Errorf("expected output to contain unresolved_mappings=0, got: %s", out)
	}
}

//nolint:funlen // test code
func TestPreviewCommand_WithRows(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	row := &commonv1.ResultEntry{}
	row.SetId(1)
	row.SetRaceGridId(2)
	row.SetDriverId(3)
	row.SetCarModelId(4)
	row.SetFinishingPosition(1)
	row.SetCompletedLaps(50)

	unresolved := &commonv1.UnresolvedMapping{}
	unresolved.SetSourceValue("driver_abc")
	unresolved.SetMappingType("driver")

	runner := &previewCommand{
		raceGridID: 2,
		out:        &buf,
		qrySvc:     &mockQueryClient{},
		importSvc: &mockImportClient{
			getPreprocessPreview: func(_ context.Context, _ *connect.Request[importv1.GetPreprocessPreviewRequest]) (*connect.Response[importv1.GetPreprocessPreviewResponse], error) {
				resp := &importv1.GetPreprocessPreviewResponse{}
				resp.SetRows([]*commonv1.ResultEntry{row})
				resp.SetUnresolvedMappings([]*commonv1.UnresolvedMapping{unresolved})
				return connect.NewResponse(resp), nil
			},
		},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "rows=1") {
		t.Errorf("expected output to contain rows=1, got: %s", out)
	}
	if !strings.Contains(out, "unresolved_mappings=1") {
		t.Errorf("expected output to contain unresolved_mappings=1, got: %s", out)
	}
	if !strings.Contains(out, "driver_id=3") {
		t.Errorf("expected output to contain driver_id=3, got: %s", out)
	}
	if !strings.Contains(out, "source_value=driver_abc") {
		t.Errorf("expected output to contain source_value=driver_abc, got: %s", out)
	}
}

func TestPreviewCommand_PreviewError(t *testing.T) {
	t.Parallel()

	runner := &previewCommand{
		raceGridID: 1,
		out:        &bytes.Buffer{},
		qrySvc:     &mockQueryClient{},
		importSvc: &mockImportClient{
			getPreprocessPreview: func(_ context.Context, _ *connect.Request[importv1.GetPreprocessPreviewRequest]) (*connect.Response[importv1.GetPreprocessPreviewResponse], error) {
				return nil, errors.New("preview failed")
			},
		},
	}

	err := runner.run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get preprocess preview") {
		t.Errorf("expected error to mention 'get preprocess preview', got: %v", err)
	}
}
