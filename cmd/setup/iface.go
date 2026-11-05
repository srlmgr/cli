package setup

import (
	"context"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
)

// commandClient is the subset of CommandServiceClient used by the setup runner.
type commandClient interface {
	CreateSimulation(
		context.Context, *commandv1.CreateSimulationRequest,
	) (*commandv1.CreateSimulationResponse, error)

	CreateSeries(
		context.Context, *commandv1.CreateSeriesRequest,
	) (*commandv1.CreateSeriesResponse, error)

	CreateSeason(
		context.Context, *commandv1.CreateSeasonRequest,
	) (*commandv1.CreateSeasonResponse, error)
	AssignCarClassToSeason(
		context.Context, *commandv1.AssignCarClassToSeasonRequest,
	) (*commandv1.AssignCarClassToSeasonResponse, error)

	CreatePointSystem(
		context.Context, *commandv1.CreatePointSystemRequest,
	) (*commandv1.CreatePointSystemResponse, error)

	CreateTrack(
		context.Context, *commandv1.CreateTrackRequest,
	) (*commandv1.CreateTrackResponse, error)

	CreateTrackLayout(
		context.Context, *commandv1.CreateTrackLayoutRequest,
	) (*commandv1.CreateTrackLayoutResponse, error)

	CreateCarManufacturer(
		context.Context, *commandv1.CreateCarManufacturerRequest,
	) (*commandv1.CreateCarManufacturerResponse, error)

	CreateCarModel(
		context.Context, *commandv1.CreateCarModelRequest,
	) (*commandv1.CreateCarModelResponse, error)
	CreateCarModelVariant(
		context.Context, *commandv1.CreateCarModelVariantRequest,
	) (*commandv1.CreateCarModelVariantResponse, error)
	CreateCarClass(
		context.Context, *commandv1.CreateCarClassRequest,
	) (*commandv1.CreateCarClassResponse, error)
	AssignCarModelVariantToCarClass(
		context.Context, *commandv1.AssignCarModelVariantToCarClassRequest,
	) (*commandv1.AssignCarModelVariantToCarClassResponse, error)

	CreateDriver(
		context.Context, *commandv1.CreateDriverRequest,
	) (*commandv1.CreateDriverResponse, error)
	CreateTeam(
		context.Context, *commandv1.CreateTeamRequest,
	) (*commandv1.CreateTeamResponse, error)
	AddSeasonDriver(
		context.Context, *commandv1.AddSeasonDriverRequest,
	) (*commandv1.AddSeasonDriverResponse, error)
	SetSeasonDrivers(
		context.Context, *commandv1.SetSeasonDriversRequest,
	) (*commandv1.SetSeasonDriversResponse, error)
	SetSeasonCarClasses(
		context.Context, *commandv1.SetSeasonCarClassesRequest,
	) (*commandv1.SetSeasonCarClassesResponse, error)
	SetSeasonCarModelVariants(
		context.Context, *commandv1.SetSeasonCarModelVariantsRequest,
	) (*commandv1.SetSeasonCarModelVariantsResponse, error)

	CreateEvent(
		context.Context, *commandv1.CreateEventRequest,
	) (*commandv1.CreateEventResponse, error)

	CreateRace(
		context.Context, *commandv1.CreateRaceRequest,
	) (*commandv1.CreateRaceResponse, error)

	CreateRaceGrid(
		context.Context, *commandv1.CreateRaceGridRequest,
	) (*commandv1.CreateRaceGridResponse, error)

	SetTeamMembers(
		context.Context, *commandv1.SetTeamMembersRequest,
	) (*commandv1.SetTeamMembersResponse, error)

	SetSimulationDriverAliases(
		context.Context, *commandv1.SetSimulationDriverAliasesRequest,
	) (*commandv1.SetSimulationDriverAliasesResponse, error)

	SetSimulationCarAliases(
		context.Context, *commandv1.SetSimulationCarAliasesRequest,
	) (*commandv1.SetSimulationCarAliasesResponse, error)

	SetSimulationTrackLayoutAliases(
		context.Context, *commandv1.SetSimulationTrackLayoutAliasesRequest,
	) (*commandv1.SetSimulationTrackLayoutAliasesResponse, error)
}

// queryClient is the subset of QueryServiceClient used by the setup runner.
type queryClient interface {
	ListSimulations(
		context.Context, *queryv1.ListSimulationsRequest,
	) (*queryv1.ListSimulationsResponse, error)

	ListSeries(
		context.Context, *queryv1.ListSeriesRequest,
	) (*queryv1.ListSeriesResponse, error)

	ListSeasons(
		context.Context, *queryv1.ListSeasonsRequest,
	) (*queryv1.ListSeasonsResponse, error)

	ListPointSystems(
		context.Context, *queryv1.ListPointSystemsRequest,
	) (*queryv1.ListPointSystemsResponse, error)

	ListTracks(
		context.Context, *queryv1.ListTracksRequest,
	) (*queryv1.ListTracksResponse, error)

	ListTrackLayouts(
		context.Context, *queryv1.ListTrackLayoutsRequest,
	) (*queryv1.ListTrackLayoutsResponse, error)

	ListCarManufacturers(
		context.Context, *queryv1.ListCarManufacturersRequest,
	) (*queryv1.ListCarManufacturersResponse, error)

	ListCarModels(
		context.Context, *queryv1.ListCarModelsRequest,
	) (*queryv1.ListCarModelsResponse, error)
	ListCarModelVariants(
		context.Context, *queryv1.ListCarModelVariantsRequest,
	) (*queryv1.ListCarModelVariantsResponse, error)
	ListCarClasses(
		context.Context, *queryv1.ListCarClassesRequest,
	) (*queryv1.ListCarClassesResponse, error)

	ListDrivers(
		context.Context, *queryv1.ListDriversRequest,
	) (*queryv1.ListDriversResponse, error)
	ListTeams(
		context.Context, *queryv1.ListTeamsRequest,
	) (*queryv1.ListTeamsResponse, error)

	ListEvents(
		context.Context, *queryv1.ListEventsRequest,
	) (*queryv1.ListEventsResponse, error)

	ListRaces(
		context.Context, *queryv1.ListRacesRequest,
	) (*queryv1.ListRacesResponse, error)

	ListRaceGrids(
		context.Context, *queryv1.ListRaceGridsRequest,
	) (*queryv1.ListRaceGridsResponse, error)
}
