//nolint:whitespace,dupl,funlen // test helper functions have similar structure
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

	createDriver func(
		context.Context,
		*connect.Request[commandv1.CreateDriverRequest],
	) (*connect.Response[commandv1.CreateDriverResponse], error)

	createEvent func(
		context.Context,
		*connect.Request[commandv1.CreateEventRequest],
	) (*connect.Response[commandv1.CreateEventResponse], error)

	createRace func(
		context.Context,
		*connect.Request[commandv1.CreateRaceRequest],
	) (*connect.Response[commandv1.CreateRaceResponse], error)
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

func (m *mockCommandClient) CreateDriver(
	ctx context.Context,
	req *connect.Request[commandv1.CreateDriverRequest],
) (*connect.Response[commandv1.CreateDriverResponse], error) {
	if m.createDriver != nil {
		return m.createDriver(ctx, req)
	}

	resp := &commandv1.CreateDriverResponse{}
	d := &commonv1.Driver{}
	d.SetName(req.Msg.GetName())
	d.SetId(40)
	resp.SetDriver(d)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateEvent(
	ctx context.Context,
	req *connect.Request[commandv1.CreateEventRequest],
) (*connect.Response[commandv1.CreateEventResponse], error) {
	if m.createEvent != nil {
		return m.createEvent(ctx, req)
	}

	resp := &commandv1.CreateEventResponse{}
	e := &commonv1.Event{}
	e.SetName(req.Msg.GetName())
	e.SetId(50)
	resp.SetEvent(e)

	return connect.NewResponse(resp), nil
}

func (m *mockCommandClient) CreateRace(
	ctx context.Context,
	req *connect.Request[commandv1.CreateRaceRequest],
) (*connect.Response[commandv1.CreateRaceResponse], error) {
	if m.createRace != nil {
		return m.createRace(ctx, req)
	}

	resp := &commandv1.CreateRaceResponse{}
	rc := &commonv1.Race{}
	rc.SetName(req.Msg.GetName())
	rc.SetId(51)
	resp.SetRace(rc)

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

	listDrivers func(
		context.Context,
		*connect.Request[queryv1.ListDriversRequest],
	) (*connect.Response[queryv1.ListDriversResponse], error)

	listEvents func(
		context.Context,
		*connect.Request[queryv1.ListEventsRequest],
	) (*connect.Response[queryv1.ListEventsResponse], error)

	listRaces func(
		context.Context,
		*connect.Request[queryv1.ListRacesRequest],
	) (*connect.Response[queryv1.ListRacesResponse], error)
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

func (m *mockQueryClient) ListDrivers(
	ctx context.Context,
	req *connect.Request[queryv1.ListDriversRequest],
) (*connect.Response[queryv1.ListDriversResponse], error) {
	if m.listDrivers != nil {
		return m.listDrivers(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListDriversResponse{}), nil
}

func (m *mockQueryClient) ListEvents(
	ctx context.Context,
	req *connect.Request[queryv1.ListEventsRequest],
) (*connect.Response[queryv1.ListEventsResponse], error) {
	if m.listEvents != nil {
		return m.listEvents(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListEventsResponse{}), nil
}

func (m *mockQueryClient) ListRaces(
	ctx context.Context,
	req *connect.Request[queryv1.ListRacesRequest],
) (*connect.Response[queryv1.ListRacesResponse], error) {
	if m.listRaces != nil {
		return m.listRaces(ctx, req)
	}

	return connect.NewResponse(&queryv1.ListRacesResponse{}), nil
}

// ---- YAML parsing tests ----

func TestLoadConfig_ValidYAML(t *testing.T) {
	t.Parallel()

	cfg, err := loadConfig("testdata/fixture.yml")
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}

	if len(cfg.Drivers) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(cfg.Drivers))
	}

	if cfg.Drivers[0].Name != "Max Verstappen" {
		t.Errorf(
			"expected driver name %q, got %q",
			"Max Verstappen", cfg.Drivers[0].Name,
		)
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

	if !sim.IsActive {
		t.Error("expected simulation isActive=true")
	}

	if len(sim.Series) != 1 || sim.Series[0].Name != "Porsche Cup" {
		t.Errorf("unexpected series in simulation: %+v", sim.Series)
	}

	seasons := sim.Series[0].Seasons
	if len(seasons) != 1 || seasons[0].Name != "Saison XVIII" {
		t.Errorf("unexpected seasons: %+v", seasons)
	}

	events := seasons[0].Events
	if len(events) != 1 || events[0].Name != "Round 1 - Interlagos" {
		t.Errorf("unexpected events: %+v", events)
	}

	if events[0].TrackLayout != "Grand Prix" {
		t.Errorf("expected trackLayout %q, got %q", "Grand Prix", events[0].TrackLayout)
	}

	races := events[0].Races
	if len(races) != 1 || races[0].Name != "Race 1" {
		t.Errorf("unexpected races: %+v", races)
	}

	if races[0].SessionType != "RACE_SESSION_TYPE_RACE" {
		t.Errorf("expected sessionType %q, got %q", "RACE_SESSION_TYPE_RACE", races[0].SessionType)
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
	assertContains(t, out, `created driver "Max Verstappen"`)
	assertContains(t, out, `created point-system "Standard"`)
	assertContains(t, out, `created simulation "iRacing"`)
	assertContains(t, out, `created series "Porsche Cup"`)
	assertContains(t, out, `created season "Saison XVIII"`)
	assertContains(t, out, `created event "Round 1 - Interlagos"`)
	assertContains(t, out, `created race "Race 1"`)
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
	assertContains(t, out, `existing driver "Max Verstappen"`)
	assertContains(t, out, `existing point-system "Standard"`)
	assertContains(t, out, `existing simulation "iRacing"`)
	assertContains(t, out, `existing series "Porsche Cup"`)
	assertContains(t, out, `existing season "Saison XVIII"`)
	assertContains(t, out, `existing event "Round 1 - Interlagos"`)
	assertContains(t, out, `existing race "Race 1"`)
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

// ---- runner: driver create / skip-if-existing ----

func TestSetupRunner_CreateDriver(t *testing.T) {
	t.Parallel()

	var capturedName string

	var capturedExternalID string

	var capturedIsActive bool

	cmd := &mockCommandClient{
		createDriver: func(
			_ context.Context,
			req *connect.Request[commandv1.CreateDriverRequest],
		) (*connect.Response[commandv1.CreateDriverResponse], error) {
			capturedName = req.Msg.GetName()
			capturedExternalID = req.Msg.GetExternalId()
			capturedIsActive = req.Msg.GetIsActive()

			resp := &commandv1.CreateDriverResponse{}
			d := &commonv1.Driver{}
			d.SetId(40)
			d.SetName(req.Msg.GetName())
			resp.SetDriver(d)

			return connect.NewResponse(resp), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedName != "Max Verstappen" {
		t.Errorf("expected driver name %q, got %q", "Max Verstappen", capturedName)
	}

	if capturedExternalID != "830" {
		t.Errorf("expected externalId=%q, got %q", "830", capturedExternalID)
	}

	if !capturedIsActive {
		t.Error("expected isActive=true")
	}

	assertContains(t, buf.String(), `created driver "Max Verstappen"`)
}

func TestSetupRunner_ExistingDriver(t *testing.T) {
	t.Parallel()

	createCalled := false
	cmd := &mockCommandClient{
		createDriver: func(
			_ context.Context,
			_ *connect.Request[commandv1.CreateDriverRequest],
		) (*connect.Response[commandv1.CreateDriverResponse], error) {
			createCalled = true

			return connect.NewResponse(&commandv1.CreateDriverResponse{}), nil
		},
	}

	qry := &mockQueryClient{
		listDrivers: func(
			_ context.Context,
			_ *connect.Request[queryv1.ListDriversRequest],
		) (*connect.Response[queryv1.ListDriversResponse], error) {
			d := &commonv1.Driver{}
			d.SetId(40)
			d.SetName("Max Verstappen")
			resp := &queryv1.ListDriversResponse{}
			resp.SetItems([]*commonv1.Driver{d})

			return connect.NewResponse(resp), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   qry,
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if createCalled {
		t.Error("expected no createDriver call for existing driver, but it was called")
	}

	assertContains(t, buf.String(), `existing driver "Max Verstappen"`)
}

// ---- runner: event create / skip-if-existing ----

func TestSetupRunner_CreateEvent(t *testing.T) {
	t.Parallel()

	var capturedSeasonID uint32

	var capturedLayoutID uint32

	var capturedName string

	cmd := &mockCommandClient{
		createEvent: func(
			_ context.Context,
			req *connect.Request[commandv1.CreateEventRequest],
		) (*connect.Response[commandv1.CreateEventResponse], error) {
			capturedSeasonID = req.Msg.GetSeasonId()
			capturedLayoutID = req.Msg.GetTrackLayoutId()
			capturedName = req.Msg.GetName()

			resp := &commandv1.CreateEventResponse{}
			e := &commonv1.Event{}
			e.SetId(50)
			e.SetName(req.Msg.GetName())
			resp.SetEvent(e)

			return connect.NewResponse(resp), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedName != "Round 1 - Interlagos" {
		t.Errorf("expected event name %q, got %q", "Round 1 - Interlagos", capturedName)
	}

	if capturedSeasonID == 0 {
		t.Error("expected non-zero seasonId on CreateEvent request")
	}

	if capturedLayoutID == 0 {
		t.Error("expected non-zero trackLayoutId on CreateEvent request")
	}

	assertContains(t, buf.String(), `created event "Round 1 - Interlagos"`)
}

func TestSetupRunner_ExistingEventSkipsRaces(t *testing.T) {
	t.Parallel()

	raceCreateCalled := false
	cmd := &mockCommandClient{
		createRace: func(
			_ context.Context,
			_ *connect.Request[commandv1.CreateRaceRequest],
		) (*connect.Response[commandv1.CreateRaceResponse], error) {
			raceCreateCalled = true

			return connect.NewResponse(&commandv1.CreateRaceResponse{}), nil
		},
	}

	qry := &mockQueryClient{
		listRaces: func(
			_ context.Context,
			_ *connect.Request[queryv1.ListRacesRequest],
		) (*connect.Response[queryv1.ListRacesResponse], error) {
			rc := &commonv1.Race{}
			rc.SetId(51)
			rc.SetName("Race 1")
			resp := &queryv1.ListRacesResponse{}
			resp.SetItems([]*commonv1.Race{rc})

			return connect.NewResponse(resp), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   qry,
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if raceCreateCalled {
		t.Error("race must not be created when it already exists")
	}

	assertContains(t, buf.String(), `existing race "Race 1"`)
}

// ---- runner: race create when event was newly provisioned ----

func TestSetupRunner_NewEventCreatesRaces(t *testing.T) {
	t.Parallel()

	var capturedRaceName string

	var capturedSessionType commonv1.RaceSessionType

	var capturedSequenceNo int32

	cmd := &mockCommandClient{
		createRace: func(
			_ context.Context,
			req *connect.Request[commandv1.CreateRaceRequest],
		) (*connect.Response[commandv1.CreateRaceResponse], error) {
			capturedRaceName = req.Msg.GetName()
			capturedSessionType = req.Msg.GetSessionType()
			capturedSequenceNo = req.Msg.GetSequenceNo()

			resp := &commandv1.CreateRaceResponse{}
			rc := &commonv1.Race{}
			rc.SetId(51)
			rc.SetName(req.Msg.GetName())
			resp.SetRace(rc)

			return connect.NewResponse(resp), nil
		},
	}

	var buf bytes.Buffer

	runner := &setupRunner{
		filePath: "testdata/fixture.yml",
		dryRun:   false,
		out:      &buf,
		cmdSvc:   cmd,
		qrySvc:   &mockQueryClient{},
	}

	if err := runner.run(context.Background()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if capturedRaceName != "Race 1" {
		t.Errorf("expected race name %q, got %q", "Race 1", capturedRaceName)
	}

	if capturedSessionType != commonv1.RaceSessionType_RACE_SESSION_TYPE_RACE {
		t.Errorf(
			"expected sessionType %v, got %v",
			commonv1.RaceSessionType_RACE_SESSION_TYPE_RACE, capturedSessionType,
		)
	}

	if capturedSequenceNo != 1 {
		t.Errorf("expected sequenceNo=1, got %d", capturedSequenceNo)
	}

	assertContains(t, buf.String(), `created race "Race 1"`)
}

// ---- runner: validation tests for new entities ----

func TestValidateConfig_MissingDriverName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Drivers: []DriverConfig{{Name: ""}},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty driver name, got nil")
	}

	if !strings.Contains(err.Error(), "drivers[0]") {
		t.Errorf("expected error to mention drivers[0], got: %v", err)
	}
}

func TestValidateConfig_MissingEventName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Simulations: []SimulationConfig{
			{
				Name: "iRacing",
				Series: []SeriesConfig{
					{
						Name: "Porsche Cup",
						Seasons: []SeasonConfig{
							{
								Name:   "Saison XVIII",
								Events: []EventConfig{{Name: ""}},
							},
						},
					},
				},
			},
		},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty event name, got nil")
	}

	if !strings.Contains(err.Error(), "events[0]") {
		t.Errorf("expected error to mention events[0], got: %v", err)
	}
}

func TestValidateConfig_MissingRaceName(t *testing.T) {
	t.Parallel()

	cfg := &SetupConfig{
		Simulations: []SimulationConfig{
			{
				Name: "iRacing",
				Series: []SeriesConfig{
					{
						Name: "Porsche Cup",
						Seasons: []SeasonConfig{
							{
								Name: "Saison XVIII",
								Events: []EventConfig{
									{
										Name:  "Round 1",
										Races: []RaceConfig{{Name: ""}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for empty race name, got nil")
	}

	if !strings.Contains(err.Error(), "races[0]") {
		t.Errorf("expected error to mention races[0], got: %v", err)
	}
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

	// Verify all 12 entity types were provisioned.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 12 {
		t.Errorf("expected 12 output lines, got %d:\n%s", len(lines), buf.String())
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
		simID:    1,
		srID:     2,
		snID:     3,
		psID:     10,
		trackID:  20,
		layID:    21,
		mfrID:    30,
		brandID:  31,
		modelID:  32,
		driverID: 40,
		eventID:  50,
		raceID:   51,
	})
}

type fixtureIDs struct {
	simID, srID, snID       uint32
	psID                    uint32
	trackID, layID          uint32
	mfrID, brandID, modelID uint32
	driverID                uint32
	eventID                 uint32
	raceID                  uint32
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
	qry.listDrivers = buildExistingDrivers(ids)
	qry.listEvents = buildExistingEvents(ids)
	qry.listRaces = buildExistingRaces(ids)

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

func buildExistingDrivers(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListDriversRequest],
) (*connect.Response[queryv1.ListDriversResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListDriversRequest],
	) (*connect.Response[queryv1.ListDriversResponse], error) {
		d := &commonv1.Driver{}
		d.SetId(ids.driverID)
		d.SetName("Max Verstappen")
		resp := &queryv1.ListDriversResponse{}
		resp.SetItems([]*commonv1.Driver{d})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingEvents(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListEventsRequest],
) (*connect.Response[queryv1.ListEventsResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListEventsRequest],
	) (*connect.Response[queryv1.ListEventsResponse], error) {
		e := &commonv1.Event{}
		e.SetId(ids.eventID)
		e.SetName("Round 1 - Interlagos")
		e.SetSeasonId(ids.snID)
		resp := &queryv1.ListEventsResponse{}
		resp.SetItems([]*commonv1.Event{e})

		return connect.NewResponse(resp), nil
	}
}

func buildExistingRaces(ids fixtureIDs) func(
	context.Context, *connect.Request[queryv1.ListRacesRequest],
) (*connect.Response[queryv1.ListRacesResponse], error) {
	return func(
		_ context.Context,
		_ *connect.Request[queryv1.ListRacesRequest],
	) (*connect.Response[queryv1.ListRacesResponse], error) {
		rc := &commonv1.Race{}
		rc.SetId(ids.raceID)
		rc.SetName("Race 1")
		rc.SetEventId(ids.eventID)
		resp := &queryv1.ListRacesResponse{}
		resp.SetItems([]*commonv1.Race{rc})

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