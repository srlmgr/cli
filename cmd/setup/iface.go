package setup

import (
	"context"

	commandv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/command/v1"
	queryv1 "buf.build/gen/go/srlmgr/api/protocolbuffers/go/backend/query/v1"
	"connectrpc.com/connect"
)

// commandClient is the subset of CommandServiceClient used by the setup runner.
type commandClient interface {
	CreateSimulation(
		context.Context,
		*connect.Request[commandv1.CreateSimulationRequest],
	) (*connect.Response[commandv1.CreateSimulationResponse], error)

	CreateSeries(
		context.Context,
		*connect.Request[commandv1.CreateSeriesRequest],
	) (*connect.Response[commandv1.CreateSeriesResponse], error)

	CreateSeason(
		context.Context,
		*connect.Request[commandv1.CreateSeasonRequest],
	) (*connect.Response[commandv1.CreateSeasonResponse], error)
	AssignCarClassToSeason(
		context.Context,
		*connect.Request[commandv1.AssignCarClassToSeasonRequest],
	) (*connect.Response[commandv1.AssignCarClassToSeasonResponse], error)

	CreatePointSystem(
		context.Context,
		*connect.Request[commandv1.CreatePointSystemRequest],
	) (*connect.Response[commandv1.CreatePointSystemResponse], error)

	CreateTrack(
		context.Context,
		*connect.Request[commandv1.CreateTrackRequest],
	) (*connect.Response[commandv1.CreateTrackResponse], error)

	CreateTrackLayout(
		context.Context,
		*connect.Request[commandv1.CreateTrackLayoutRequest],
	) (*connect.Response[commandv1.CreateTrackLayoutResponse], error)

	CreateCarManufacturer(
		context.Context,
		*connect.Request[commandv1.CreateCarManufacturerRequest],
	) (*connect.Response[commandv1.CreateCarManufacturerResponse], error)

	CreateCarBrand(
		context.Context,
		*connect.Request[commandv1.CreateCarBrandRequest],
	) (*connect.Response[commandv1.CreateCarBrandResponse], error)

	CreateCarModel(
		context.Context,
		*connect.Request[commandv1.CreateCarModelRequest],
	) (*connect.Response[commandv1.CreateCarModelResponse], error)
	CreateCarModelV2(
		context.Context,
		*connect.Request[commandv1.CreateCarModelV2Request],
	) (*connect.Response[commandv1.CreateCarModelV2Response], error)
	CreateCarModelVariant(
		context.Context,
		*connect.Request[commandv1.CreateCarModelVariantRequest],
	) (*connect.Response[commandv1.CreateCarModelVariantResponse], error)
	CreateCarClass(
		context.Context,
		*connect.Request[commandv1.CreateCarClassRequest],
	) (*connect.Response[commandv1.CreateCarClassResponse], error)
	AssignCarModelToCarClass(
		context.Context,
		*connect.Request[commandv1.AssignCarModelToCarClassRequest],
	) (*connect.Response[commandv1.AssignCarModelToCarClassResponse], error)

	CreateDriver(
		context.Context,
		*connect.Request[commandv1.CreateDriverRequest],
	) (*connect.Response[commandv1.CreateDriverResponse], error)
	CreateTeam(
		context.Context,
		*connect.Request[commandv1.CreateTeamRequest],
	) (*connect.Response[commandv1.CreateTeamResponse], error)
	AddSeasonDriver(
		context.Context,
		*connect.Request[commandv1.AddSeasonDriverRequest],
	) (*connect.Response[commandv1.AddSeasonDriverResponse], error)
	SetSeasonDrivers(
		context.Context,
		*connect.Request[commandv1.SetSeasonDriversRequest],
	) (*connect.Response[commandv1.SetSeasonDriversResponse], error)
	SetSeasonCarClasses(
		context.Context,
		*connect.Request[commandv1.SetSeasonCarClassesRequest],
	) (*connect.Response[commandv1.SetSeasonCarClassesResponse], error)
	SetSeasonCarModels(
		context.Context,
		*connect.Request[commandv1.SetSeasonCarModelsRequest],
	) (*connect.Response[commandv1.SetSeasonCarModelsResponse], error)

	CreateEvent(
		context.Context,
		*connect.Request[commandv1.CreateEventRequest],
	) (*connect.Response[commandv1.CreateEventResponse], error)

	CreateRace(
		context.Context,
		*connect.Request[commandv1.CreateRaceRequest],
	) (*connect.Response[commandv1.CreateRaceResponse], error)

	CreateRaceGrid(
		context.Context,
		*connect.Request[commandv1.CreateRaceGridRequest],
	) (*connect.Response[commandv1.CreateRaceGridResponse], error)

	SetTeamMembers(
		context.Context,
		*connect.Request[commandv1.SetTeamMembersRequest],
	) (*connect.Response[commandv1.SetTeamMembersResponse], error)

	SetSimulationDriverAliases(
		context.Context,
		*connect.Request[commandv1.SetSimulationDriverAliasesRequest],
	) (*connect.Response[commandv1.SetSimulationDriverAliasesResponse], error)

	SetSimulationCarAliases(
		context.Context,
		*connect.Request[commandv1.SetSimulationCarAliasesRequest],
	) (*connect.Response[commandv1.SetSimulationCarAliasesResponse], error)

	SetSimulationTrackLayoutAliases(
		context.Context,
		*connect.Request[commandv1.SetSimulationTrackLayoutAliasesRequest],
	) (*connect.Response[commandv1.SetSimulationTrackLayoutAliasesResponse], error)
}

// queryClient is the subset of QueryServiceClient used by the setup runner.
type queryClient interface {
	ListSimulations(
		context.Context,
		*connect.Request[queryv1.ListSimulationsRequest],
	) (*connect.Response[queryv1.ListSimulationsResponse], error)

	ListSeries(
		context.Context,
		*connect.Request[queryv1.ListSeriesRequest],
	) (*connect.Response[queryv1.ListSeriesResponse], error)

	ListSeasons(
		context.Context,
		*connect.Request[queryv1.ListSeasonsRequest],
	) (*connect.Response[queryv1.ListSeasonsResponse], error)

	ListPointSystems(
		context.Context,
		*connect.Request[queryv1.ListPointSystemsRequest],
	) (*connect.Response[queryv1.ListPointSystemsResponse], error)

	ListTracks(
		context.Context,
		*connect.Request[queryv1.ListTracksRequest],
	) (*connect.Response[queryv1.ListTracksResponse], error)

	ListTrackLayouts(
		context.Context,
		*connect.Request[queryv1.ListTrackLayoutsRequest],
	) (*connect.Response[queryv1.ListTrackLayoutsResponse], error)

	ListCarManufacturers(
		context.Context,
		*connect.Request[queryv1.ListCarManufacturersRequest],
	) (*connect.Response[queryv1.ListCarManufacturersResponse], error)

	ListCarBrands(
		context.Context,
		*connect.Request[queryv1.ListCarBrandsRequest],
	) (*connect.Response[queryv1.ListCarBrandsResponse], error)

	ListCarModels(
		context.Context,
		*connect.Request[queryv1.ListCarModelsRequest],
	) (*connect.Response[queryv1.ListCarModelsResponse], error)
	ListCarModelsV2(
		context.Context,
		*connect.Request[queryv1.ListCarModelsV2Request],
	) (*connect.Response[queryv1.ListCarModelsV2Response], error)
	ListCarModelVariants(
		context.Context,
		*connect.Request[queryv1.ListCarModelVariantsRequest],
	) (*connect.Response[queryv1.ListCarModelVariantsResponse], error)
	ListCarClasses(
		context.Context,
		*connect.Request[queryv1.ListCarClassesRequest],
	) (*connect.Response[queryv1.ListCarClassesResponse], error)

	ListDrivers(
		context.Context,
		*connect.Request[queryv1.ListDriversRequest],
	) (*connect.Response[queryv1.ListDriversResponse], error)
	ListTeams(
		context.Context,
		*connect.Request[queryv1.ListTeamsRequest],
	) (*connect.Response[queryv1.ListTeamsResponse], error)

	ListEvents(
		context.Context,
		*connect.Request[queryv1.ListEventsRequest],
	) (*connect.Response[queryv1.ListEventsResponse], error)

	ListRaces(
		context.Context,
		*connect.Request[queryv1.ListRacesRequest],
	) (*connect.Response[queryv1.ListRacesResponse], error)

	ListRaceGrids(
		context.Context,
		*connect.Request[queryv1.ListRaceGridsRequest],
	) (*connect.Response[queryv1.ListRaceGridsResponse], error)
}
