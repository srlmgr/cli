package setup

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	commonv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/common/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
)

// ---- mock implementations ----

type mockCommandClient struct {
	createSimulation func(
		context.Context,
		*connect.Request[commandv1.CreateSimulationRequest],
	) (*connect.Response[commandv1.CreateSimulationResponse], error)

	createSeries func(
		context.Context,
		*connect.Request[commandv1.CreateSeriesRequest],
	) (*connect.Response[commandv1.CreateSeriesResponse], error)

	createSeason func(
		context.Context,
		*connect.Request[commandv1.CreateSeasonRequest],
	) (*connect.Response[commandv1.CreateSeasonResponse], error)

	createPointSystem func(
		context.Context,
		*connect.Request[commandv1.CreatePointSystemRequest],
	) (*connect.Response[commandv1.CreatePointSystemResponse], error)

	createTrack func(
		context.Context,
		*connect.Request[commandv1.CreateTrackRequest],
	) (*connect.Response[commandv1.CreateTrackResponse], error)

	createTrackLayout func(
		context.Context,
		*connect.Request[commandv1.CreateTrackLayoutRequest],
	) (*connect.Response[commandv1.CreateTrackLayoutResponse], error)

	createCarManufacturer func(
		context.Context,
		*connect.Request[commandv1.CreateCarManufacturerRequest],
	) (*connect.Response[commandv1.CreateCarManufacturerResponse], error)

	createCarBrand func(
		context.Context,
		*connect.Request[commandv1.CreateCarBrandRequest],
	) (*connect.Response[commandv1.CreateCarBrandResponse], error)

	createCarModel func(
		context.Context,
		*connect.Request[commandv1.CreateCarModelRequest],
	) (*connect.Response[commandv1.CreateCarModelResponse], error)
}

func (m *mockCommandClient) CreateSimulation(
	ctx context.Context,
	req *connect.Request[commandv1.CreateSimulationRequest],
) (*connect.Response[commandv1.CreateSimulationResponse], error) {
	if m.createSimulation != nil {
		return m.createSimulation(ctx, req)
	}

	resp := &commandv1.CreateSimulationResponse{}
	sim := &commonv1.Simulation{}
	sim.SetName(req.Msg.GetName())
	sim.SetId(1)
	resp.SetSimulation(sim)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateSeries(
	ctx context.Context,
	req *connect.Request[commandv1.CreateSeriesRequest],
) (*connect.Response[commandv1.CreateSeriesResponse], error) {
	if m.createSeries != nil {
		return m.createSeries(ctx, req)
	}

	resp := &commandv1.CreateSeriesResponse{}
	s := &commonv1.Series{}
	s.SetName(req.Msg.GetName())
	s.SetId(2)
	resp.SetSeries(s)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateSeason(
	ctx context.Context,
	req *connect.Request[commandv1.CreateSeasonRequest],
) (*connect.Response[commandv1.CreateSeasonResponse], error) {
	if m.createSeason != nil {
		return m.createSeason(ctx, req)
	}

	resp := &commandv1.CreateSeasonResponse{}
	s := &commonv1.Season{}
	s.SetName(req.Msg.GetName())
	s.SetId(3)
	resp.SetSeason(s)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreatePointSystem(
	ctx context.Context,
	req *connect.Request[commandv1.CreatePointSystemRequest],
) (*connect.Response[commandv1.CreatePointSystemResponse], error) {
	if m.createPointSystem != nil {
		return m.createPointSystem(ctx, req)
	}

	resp := &commandv1.CreatePointSystemResponse{}
	ps := &commonv1.PointSystem{}
	ps.SetName(req.Msg.GetName())
	ps.SetId(10)
	resp.SetPointSystem(ps)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateTrack(
	ctx context.Context,
	req *connect.Request[commandv1.CreateTrackRequest],
) (*connect.Response[commandv1.CreateTrackResponse], error) {
	if m.createTrack != nil {
		return m.createTrack(ctx, req)
	}

	resp := &commandv1.CreateTrackResponse{}
	t := &commonv1.Track{}
	t.SetName(req.Msg.GetName())
	t.SetId(20)
	resp.SetTrack(t)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateTrackLayout(
	ctx context.Context,
	req *connect.Request[commandv1.CreateTrackLayoutRequest],
) (*connect.Response[commandv1.CreateTrackLayoutResponse], error) {
	if m.createTrackLayout != nil {
		return m.createTrackLayout(ctx, req)
	}

	resp := &commandv1.CreateTrackLayoutResponse{}
	l := &commonv1.TrackLayout{}
	l.SetName(req.Msg.GetName())
	l.SetId(21)
	resp.SetTrackLayout(l)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateCarManufacturer(
	ctx context.Context,
	req *connect.Request[commandv1.CreateCarManufacturerRequest],
) (*connect.Response[commandv1.CreateCarManufacturerResponse], error) {
	if m.createCarManufacturer != nil {
		return m.createCarManufacturer(ctx, req)
	}

	resp := &commandv1.CreateCarManufacturerResponse{}
	mfr := &commonv1.CarManufacturer{}
	mfr.SetName(req.Msg.GetName())
	mfr.SetId(30)
	resp.SetCarManufacturer(mfr)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateCarBrand(
	ctx context.Context,
	req *connect.Request[commandv1.CreateCarBrandRequest],
) (*connect.Response[commandv1.CreateCarBrandResponse], error) {
	if m.createCarBrand != nil {
		return m.createCarBrand(ctx, req)
	}

	resp := &commandv1.CreateCarBrandResponse{}
	b := &commonv1.CarBrand{}
	b.SetName(req.Msg.GetName())
	b.SetId(31)
	resp.SetCarBrand(b)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateCarModel(
	ctx context.Context,
	req *connect.Request[commandv1.CreateCarModelRequest],
) (*connect.Response[commandv1.CreateCarModelResponse], error) {
	if m.createCarModel != nil {
		return m.createCarModel(ctx, req)
	}

	resp := &commandv1.CreateCarModelResponse{}
	cm := &commonv1.CarModel{}
	cm.SetName(req.Msg.GetName())
	cm.SetId(32)
	resp.SetCarModel(cm)

	return connect.NewResponse(resp), nil
}

// mockQueryClient returns empty lists by default (simulates no existing entities).
type mockQueryClient struct {
	listSimulations func(
		context.Context,
		*connect.Request[queryv1.ListSimulationsRequest],
	) (*connect.Response[queryv1.ListSimulationsResponse], error)

	listSeries func(
		context.Context,
		*connect.Request[queryv1.ListSeriesRequest],
	) (*connect.Response[queryv1.ListSeriesResponse], error)

	listSeasons func(
		context.Context,
		*connect.Request[queryv1.ListSeasonsRequest],
	) (*connect.Response[queryv1.ListSeasonsResponse], error)

	listPointSystems func(
		context.Context,
		*connect.Request[queryv1.ListPointSystemsRequest],
	) (*connect.Response[queryv1.ListPointSystemsResponse], error)

	listTracks func(
		context.Context,
		*connect.Request[queryv1.ListTracksRequest],
	) (*connect.Response[queryv1.ListTracksResponse], error)

	listTrackLayouts func(
		context.Context,
		*connect.Request[queryv1.ListTrackLayoutsRequest],
	) (*connect.Response[queryv1.ListTrackLayoutsResponse], error)

	listCarManufacturers func(
		context.Context,
		*connect.Request[queryv1.ListCarManufacturersRequest],
	) (*connect.Response[queryv1.ListCarManufacturersResponse], error)

	listCarBrands func(
		context.Context,
		*connect.Request[queryv1.ListCarBrandsRequest],
	) (*connect.Response[queryv1.ListCarBrandsResponse], error)

	listCarModels func(
		context.Context,
		*connect.Request[queryv1.ListCarModelsRequest],
	) (*connect.Response[queryv1.ListCarModelsResponse], error)
}

func (m *mockQueryClient) ListSimulations(
	ctx context.Context,
	req *connect.Request[queryv1.ListSimulationsRequest],
) (*connect.Response[queryv1.ListSimulationsResponse], error) {
	if m.listSimulations != nil {
		return m.listSimulations(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListSimulationsResponse{}), nil
}

func (m *mockQueryClient) ListSeries(
	ctx context.Context,
	req *connect.Request[queryv1.ListSeriesRequest],
) (*connect.Response[queryv1.ListSeriesResponse], error) {
	if m.listSeries != nil {
		return m.listSeries(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListSeriesResponse{}), nil
}

func (m *mockQueryClient) ListSeasons(
	ctx context.Context,
	req *connect.Request[queryv1.ListSeasonsRequest],
) (*connect.Response[queryv1.ListSeasonsResponse], error) {
	if m.listSeasons != nil {
		return m.listSeasons(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListSeasonsResponse{}), nil
}

func (m *mockQueryClient) ListPointSystems(
	ctx context.Context,
	req *connect.Request[queryv1.ListPointSystemsRequest],
) (*connect.Response[queryv1.ListPointSystemsResponse], error) {
	if m.listPointSystems != nil {
		return m.listPointSystems(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListPointSystemsResponse{}), nil
}

func (m *mockQueryClient) ListTracks(
	ctx context.Context,
	req *connect.Request[queryv1.ListTracksRequest],
) (*connect.Response[queryv1.ListTracksResponse], error) {
	if m.listTracks != nil {
		return m.listTracks(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListTracksResponse{}), nil
}

func (m *mockQueryClient) ListTrackLayouts(
	ctx context.Context,
	req *connect.Request[queryv1.ListTrackLayoutsRequest],
) (*connect.Response[queryv1.ListTrackLayoutsResponse], error) {
	if m.listTrackLayouts != nil {
		return m.listTrackLayouts(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListTrackLayoutsResponse{}), nil
}

func (m *mockQueryClient) ListCarManufacturers(
	ctx context.Context,
	req *connect.Request[queryv1.ListCarManufacturersRequest],
) (*connect.Response[queryv1.ListCarManufacturersResponse], error) {
	if m.listCarManufacturers != nil {
		return m.listCarManufacturers(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListCarManufacturersResponse{}), nil
}

func (m *mockQueryClient) ListCarBrands(
	ctx context.Context,
	req *connect.Request[queryv1.ListCarBrandsRequest],
) (*connect.Response[queryv1.ListCarBrandsResponse], error) {
	if m.listCarBrands != nil {
		return m.listCarBrands(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListCarBrandsResponse{}), nil
}

func (m *mockQueryClient) ListCarModels(
	ctx context.Context,
	req *connect.Request[queryv1.ListCarModelsRequest],
) (*connect.Response[queryv1.ListCarModelsResponse], error) {
	if m.listCarModels != nil {
		return m.listCarModels(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListCarModelsResponse{}), nil
}

// ---- YAML parsing tests ----

func TestLoadConfig_ValidYAML(t *testing.T) {
	t.Parallel()

	cfg, err := loadConfig("testdata/fixture.yml")
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}

	if len(cfg.PointSystems) != 1 {
		t.Fatalf("expected 1 point system, got %d", len(cfg.PointSystems))
	}

	if cfg.PointSystems[0].Name != "Standard" {
		t.Errorf(
			"expected point system name %q, got %q",
			"Standard", cfg.PointSystems[0].Name,
		)
	}

	if len(cfg.Simulations) != 1 {
		t.Fatalf("expected 1 simulation, got %d", len(cfg.Simulations))
	}

	sim := cfg.Simulations[0]
	if sim.Name != "iRacing" {
		t.Errorf("expected simulation name %q, got %q", "iRacing", sim.Name)
	}

	if len(sim.Series) != 1 || sim.Series[0].Name != "Porsche Cup" {
		t.Errorf("unexpected series in simulation: %+v", sim.Series)
	}

	if len(cfg.Tracks) != 1 || cfg.Tracks[0].Name != "Interlagos" {
		t.Errorf("unexpected tracks: %+v", cfg.Tracks)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := loadConfig("testdata/nonexistent.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	t.Parallel()

	f := t.TempDir() + "/bad.yml"
	if err := os.WriteFile(f, []byte("{\ninvalid: yaml: content\n"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := loadConfig(f)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// ---- validation tests ----

func TestValidateConfig_MissingPointSystemName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		PointSystems: []PointSystemConfig{{Name: ""}},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty point system name, got nil")
	}

	if !strings.Contains(err.Error(), "pointSystems[0]") {
		t.Errorf("expected error to mention pointSystems[0], got: %v", err)
	}
}

func TestValidateConfig_MissingSimulationName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Simulations: []SimulationConfig{{Name: ""}},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty simulation name, got nil")
	}

	if !strings.Contains(err.Error(), "simulations[0]") {
		t.Errorf("expected error to mention simulations[0], got: %v", err)
	}
}

func TestValidateConfig_MissingSeriesName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Simulations: []SimulationConfig{
			{
				Name:   "iRacing",
				Series: []SeriesConfig{{Name: ""}},
			},
		},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty series name, got nil")
	}

	if !strings.Contains(err.Error(), "series[0]") {
		t.Errorf("expected error to mention series[0], got: %v", err)
	}
}

func TestValidateConfig_MissingTrackName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Tracks: []TrackConfig{{Name: ""}},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty track name, got nil")
	}

	if !strings.Contains(err.Error(), "tracks[0]") {
		t.Errorf("expected error to mention tracks[0], got: %v", err)
	}
}

// ---- runner: create path ----

func TestSetupRunner_CreateEntities(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   &mockCommandClient{},
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	assertContains(t, out, `created point-system "Standard"`)
	assertContains(t, out, `created simulation "iRacing"`)
	assertContains(t, out, `created series "Porsche Cup"`)
	assertContains(t, out, `created season "Saison XVIII"`)
	assertContains(t, out, `created car-manufacturer "Porsche"`)
	assertContains(t, out, `created car-brand "Porsche 911"`)
	assertContains(t, out, `created car-model "Porsche 911 GT3 Cup (992)"`)
	assertContains(t, out, `created track "Interlagos"`)
	assertContains(t, out, `created track-layout "Grand Prix"`)
}

// ---- runner: idempotency (all entities already exist) ----

func TestSetupRunner_ExistingEntities(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   &mockCommandClient{},
		qrySvc:   existingEntitiesQueryClient(),
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	out := buf.String()
	assertContains(t, out, `existing point-system "Standard"`)
	assertContains(t, out, `existing simulation "iRacing"`)
	assertContains(t, out, `existing series "Porsche Cup"`)
	assertContains(t, out, `existing season "Saison XVIII"`)
	assertContains(t, out, `existing car-manufacturer "Porsche"`)
	assertContains(t, out, `existing car-brand "Porsche 911"`)
	assertContains(t, out, `existing car-model "Porsche 911 GT3 Cup (992)"`)
	assertContains(t, out, `existing track "Interlagos"`)
	assertContains(t, out, `existing track-layout "Grand Prix"`)
}

// ---- runner: dry-run mode ----

func TestSetupRunner_DryRun(t *testing.T) {
	t.Parallel()

	createCalled := false
	cmd := &mockCommandClient{
		createSimulation: func(
			_ context.Context,
			_ *connect.Request[commandv1.CreateSimulationRequest],
		) (*connect.Response[commandv1.CreateSimulationResponse], error) {
			createCalled = true

			return connect.NewResponse(&commandv1.CreateSimulationResponse{}), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   true,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if createCalled {
		t.Error("expected no create calls in dry-run mode, but create was called")
	}

	assertContains(t, buf.String(), `dry-run: would create simulation "iRacing"`)
}

// ---- runner: parent-child scoping ----

func TestSetupRunner_ParentChildScoping(t *testing.T) {
	t.Parallel()

	const wantSimID = uint32(42)

	var capturedSimID uint32

	qry := &mockQueryClient{
		listSeries: func(
			_ context.Context,
			req *connect.Request[queryv1.ListSeriesRequest],
		) (*connect.Response[queryv1.ListSeriesResponse], error) {
			capturedSimID = req.Msg.GetSimulationId()

			return connect.NewResponse(&queryv1.ListSeriesResponse{}), nil
		},
	}

	qry.listSimulations = func(
		_ context.Context,
		_ *connect.Request[queryv1.ListSimulationsRequest],
	) (*connect.Response[queryv1.ListSimulationsResponse], error) {
		sim := &commonv1.Simulation{}
		sim.SetId(wantSimID)
		sim.SetName("iRacing")
		resp := &queryv1.ListSimulationsResponse{}
		resp.SetItems([]*commonv1.Simulation{sim})

		return connect.NewResponse(resp), nil
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   true,
		out:      &buf,
		cmdSvc:   &mockCommandClient{},
		qrySvc:   qry,
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedSimID != wantSimID {
		t.Errorf(
			"series list scoped to simulation id=%d, want id=%d",
			capturedSimID, wantSimID,
		)
	}
}

// ---- end-to-end test using YAML fixture ----

func TestSetupRunner_E2E_FixtureFile(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   &mockCommandClient{},
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	// Verify all 9 entity types were provisioned.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 9 {
		t.Errorf("expected 9 output lines, got %d:\n%s", len(lines), buf.String())
	}

	// Running again with all entities existing must not create anything.
	buf.Reset()
	runner.qrySvc = existingEntitiesQueryClient()

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("second run returned error: %v", err)
	}

	if strings.Contains(buf.String(), "created") {
		t.Errorf("second run must not create anything; got:\n%s", buf.String())
	}
}

// ---- fixture helpers ----

// existingEntitiesQueryClient returns a mockQueryClient that reports all
// fixture entities as already existing, so the runner never calls create.
func existingEntitiesQueryClient() *mockQueryClient {
	return buildExistingQueryClient(fixtureIDs{
		simID:   1,
		srID:    2,
		snID:    3,
		psID:    10,
		trackID: 20,
		layID:   21,
		mfrID:   30,
		brandID: 31,
		modelID: 32,
	})
}

type fixtureIDs struct {
	simID, srID, snID       uint32
	psID                    uint32
	trackID, layID          uint32
	mfrID, brandID, modelID uint32
}

func buildExistingQueryClient(ids fixtureIDs) *mockQueryClient {
	qry := &mockQueryClient{}
	qry.listSimulations = buildExistingSimulations(ids)
	qry.listSeries = buildExistingSeries(ids)
	qry.listSeasons = buildExistingSeasons(ids)
	qry.listPointSystems = buildExistingPointSystems(ids)
	qry.listTracks = buildExistingTracks(ids)
	qry.listTrackLayouts = buildExistingTrackLayouts(ids)
	qry.listCarManufacturers = buildExistingCarManufacturers(ids)
	qry.listCarBrands = buildExistingCarBrands(ids)
	qry.listCarModels = buildExistingCarModels(ids)

	return qry
}

func buildExistingSimulations(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListSimulationsRequest],
) (*connect.Response[queryv1.ListSimulationsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListSimulationsRequest],
	) (*connect.Response[queryv1.ListSimulationsResponse], error) {
		sim := &commonv1.Simulation{}
		sim.SetId(ids.simID)
		sim.SetName("iRacing")
		resp := &queryv1.ListSimulationsResponse{}
		resp.SetItems([]*commonv1.Simulation{sim})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingSeries(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListSeriesRequest],
) (*connect.Response[queryv1.ListSeriesResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListSeriesRequest],
	) (*connect.Response[queryv1.ListSeriesResponse], error) {
		s := &commonv1.Series{}
		s.SetId(ids.srID)
		s.SetName("Porsche Cup")
		s.SetSimulationId(ids.simID)
		resp := &queryv1.ListSeriesResponse{}
		resp.SetItems([]*commonv1.Series{s})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingSeasons(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListSeasonsRequest],
) (*connect.Response[queryv1.ListSeasonsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListSeasonsRequest],
	) (*connect.Response[queryv1.ListSeasonsResponse], error) {
		s := &commonv1.Season{}
		s.SetId(ids.snID)
		s.SetName("Saison XVIII")
		s.SetSeriesId(ids.srID)
		resp := &queryv1.ListSeasonsResponse{}
		resp.SetItems([]*commonv1.Season{s})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingPointSystems(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListPointSystemsRequest],
) (*connect.Response[queryv1.ListPointSystemsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListPointSystemsRequest],
	) (*connect.Response[queryv1.ListPointSystemsResponse], error) {
		ps := &commonv1.PointSystem{}
		ps.SetId(ids.psID)
		ps.SetName("Standard")
		resp := &queryv1.ListPointSystemsResponse{}
		resp.SetItems([]*commonv1.PointSystem{ps})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingTracks(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListTracksRequest],
) (*connect.Response[queryv1.ListTracksResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListTracksRequest],
	) (*connect.Response[queryv1.ListTracksResponse], error) {
		t := &commonv1.Track{}
		t.SetId(ids.trackID)
		t.SetName("Interlagos")
		resp := &queryv1.ListTracksResponse{}
		resp.SetItems([]*commonv1.Track{t})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingTrackLayouts(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListTrackLayoutsRequest],
) (*connect.Response[queryv1.ListTrackLayoutsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListTrackLayoutsRequest],
	) (*connect.Response[queryv1.ListTrackLayoutsResponse], error) {
		l := &commonv1.TrackLayout{}
		l.SetId(ids.layID)
		l.SetName("Grand Prix")
		l.SetTrackId(ids.trackID)
		resp := &queryv1.ListTrackLayoutsResponse{}
		resp.SetItems([]*commonv1.TrackLayout{l})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingCarManufacturers(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListCarManufacturersRequest],
) (*connect.Response[queryv1.ListCarManufacturersResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListCarManufacturersRequest],
	) (*connect.Response[queryv1.ListCarManufacturersResponse], error) {
		m := &commonv1.CarManufacturer{}
		m.SetId(ids.mfrID)
		m.SetName("Porsche")
		resp := &queryv1.ListCarManufacturersResponse{}
		resp.SetItems([]*commonv1.CarManufacturer{m})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingCarBrands(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListCarBrandsRequest],
) (*connect.Response[queryv1.ListCarBrandsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListCarBrandsRequest],
	) (*connect.Response[queryv1.ListCarBrandsResponse], error) {
		b := &commonv1.CarBrand{}
		b.SetId(ids.brandID)
		b.SetName("Porsche 911")
		b.SetManufacturerId(ids.mfrID)
		resp := &queryv1.ListCarBrandsResponse{}
		resp.SetItems([]*commonv1.CarBrand{b})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingCarModels(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListCarModelsRequest],
) (*connect.Response[queryv1.ListCarModelsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListCarModelsRequest],
	) (*connect.Response[queryv1.ListCarModelsResponse], error) {
		m := &commonv1.CarModel{}
		m.SetId(ids.modelID)
		m.SetName("Porsche 911 GT3 Cup (992)")
		m.SetBrandId(ids.brandID)
		resp := &queryv1.ListCarModelsResponse{}
		resp.SetItems([]*commonv1.CarModel{m})

		return connect.NewResponse(resp), nil
	}
}

// ---- test helpers ----

func assertContains(t *testing.T, s, substr string) {
	t.Helper()

	if !strings.Contains(s, substr) {
		t.Errorf("expected output to contain %q\ngot:\n%s", substr, s)
	}
}
