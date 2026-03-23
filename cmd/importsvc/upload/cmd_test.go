package upload

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	importv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/import/v1"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
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
	uploadResultsFile func(context.Context, *connect.Request[importv1.UploadResultsFileRequest]) (*connect.Response[importv1.UploadResultsFileResponse], error)
}

func (m *mockImportClient) UploadResultsFile(
	ctx context.Context,
	req *connect.Request[importv1.UploadResultsFileRequest],
) (*connect.Response[importv1.UploadResultsFileResponse], error) {
	if m.uploadResultsFile != nil {
		return m.uploadResultsFile(ctx, req)
	}

	resp := &importv1.UploadResultsFileResponse{}
	resp.SetImportBatchId(42)
	resp.SetProcessingState("PENDING")

	return connect.NewResponse(resp), nil
}

// ---- tests ----

func TestUploadCommand_Success(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &uploadCommand{
		raceID:       1,
		importFormat: "json",
		payload:      []byte(`{"results": []}`),
		out:          &buf,
		qrySvc:       &mockQueryClient{},
		importSvc:    &mockImportClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "import_batch_id=42") {
		t.Errorf("expected output to contain import_batch_id=42, got: %s", out)
	}
	if !strings.Contains(out, "processing_state=PENDING") {
		t.Errorf("expected output to contain processing_state=PENDING, got: %s", out)
	}
}

func TestUploadCommand_ResolvesEventID(t *testing.T) {
	t.Parallel()

	var capturedEventID uint32

	runner := &uploadCommand{
		raceID:       5,
		importFormat: "csv",
		payload:      []byte("col1,col2\nval1,val2"),
		out:          &bytes.Buffer{},
		qrySvc: &mockQueryClient{
			getRace: func(_ context.Context, req *connect.Request[queryv1.GetRaceRequest]) (*connect.Response[queryv1.GetRaceResponse], error) {
				race := &commonv1.Race{}
				race.SetId(req.Msg.GetId())
				race.SetEventId(99)
				resp := &queryv1.GetRaceResponse{}
				resp.SetRace(race)
				return connect.NewResponse(resp), nil
			},
		},
		importSvc: &mockImportClient{
			uploadResultsFile: func(_ context.Context, req *connect.Request[importv1.UploadResultsFileRequest]) (*connect.Response[importv1.UploadResultsFileResponse], error) {
				capturedEventID = req.Msg.GetEventId()
				resp := &importv1.UploadResultsFileResponse{}
				resp.SetImportBatchId(1)
				resp.SetProcessingState("OK")
				return connect.NewResponse(resp), nil
			},
		},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedEventID != 99 {
		t.Errorf("expected event_id=99, got %d", capturedEventID)
	}
}

func TestUploadCommand_GetRaceError(t *testing.T) {
	t.Parallel()

	runner := &uploadCommand{
		raceID:       1,
		importFormat: "json",
		payload:      []byte("{}"),
		out:          &bytes.Buffer{},
		qrySvc: &mockQueryClient{
			getRace: func(_ context.Context, _ *connect.Request[queryv1.GetRaceRequest]) (*connect.Response[queryv1.GetRaceResponse], error) {
				return nil, errors.New("race not found")
			},
		},
		importSvc: &mockImportClient{},
	}

	err := runner.run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get race") {
		t.Errorf("expected error to mention 'get race', got: %v", err)
	}
}

func TestUploadCommand_InvalidImportFormat(t *testing.T) {
	t.Parallel()

	runner := &uploadCommand{
		raceID:       1,
		importFormat: "invalid-format",
		payload:      []byte("{}"),
		out:          &bytes.Buffer{},
		qrySvc:       &mockQueryClient{},
		importSvc:    &mockImportClient{},
	}

	err := runner.run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse import format") {
		t.Errorf("expected error to mention 'parse import format', got: %v", err)
	}
}

func TestUploadCommand_UploadError(t *testing.T) {
	t.Parallel()

	runner := &uploadCommand{
		raceID:       1,
		importFormat: "json",
		payload:      []byte("{}"),
		out:          &bytes.Buffer{},
		qrySvc:       &mockQueryClient{},
		importSvc: &mockImportClient{
			uploadResultsFile: func(_ context.Context, _ *connect.Request[importv1.UploadResultsFileRequest]) (*connect.Response[importv1.UploadResultsFileResponse], error) {
				return nil, errors.New("upload failed")
			},
		},
	}

	err := runner.run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "upload results file") {
		t.Errorf("expected error to mention 'upload results file', got: %v", err)
	}
}
